package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/clients"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/utils"
)

type AIChatService struct {
	aiChatStorage  storage.IAIChatStorage
	tripStorage    storage.ITripStorage
	sessionStorage storage.ISessionStorage
	notifyUtils    utils.NotifyUtils
	openAIClient   clients.IChatClient
	googleApi      clients.IGoogleApiClient
}

func NewAIChatService(
	aiChatStorage storage.IAIChatStorage,
	tripStorage storage.ITripStorage,
	sessionStorage storage.ISessionStorage,
	notifyUtils utils.NotifyUtils,
	openAIClient clients.IChatClient,
	googleApi clients.IGoogleApiClient,
) service.IAIChatService {
	return &AIChatService{
		aiChatStorage:  aiChatStorage,
		tripStorage:    tripStorage,
		sessionStorage: sessionStorage,
		notifyUtils:    notifyUtils,
		openAIClient:   openAIClient,
		googleApi:      googleApi,
	}
}

func (s *AIChatService) GetAIChatMessages(ctx context.Context, tripID uuid.UUID) ([]model.ChatMessage, error) {
	messages, err := s.aiChatStorage.GetMessagesByTripID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by trip id %v: %w", tripID, err)
	}
	return messages, nil
}

func (s *AIChatService) SentMessage(ctx context.Context, message model.ChatMessage, userID int) error {
	err := s.notifyUtils.FormAndSendNotifyMessage(
		ctx,
		message.TripID,
		"chat_freeze",
		"Ваш запрос принят в обработку",
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to send chat_freeze notify: %w", err)
	}

	messageHistory, err := s.aiChatStorage.GetMessagesByTripID(ctx, message.TripID)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to get chat history for trip %v: %w", message.TripID, err)
	}

	trip, err := s.tripStorage.GetTripByID(ctx, message.TripID)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to get trip %v: %w", message.TripID, err)
	}

	prompt, err := s.formPromptOverUserMessage(message, trip.Area.GooglePlace.Name)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to form prompt over user message: %w", err)
	}

	messageHistory = append(messageHistory, model.ChatMessage{
		Role:    message.Role,
		Content: prompt,
	})

	promptResp, err := s.openAIClient.PostPrompt(ctx, messageHistory, clients.ModelChatGPT4o)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to post prompt: %w", err)
	}

	replyMessage, err := s.processPlannerResponse(ctx, promptResp, trip.Area.GooglePlace.Name)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to process planner response: %w", err)
	}

	//save user question
	err = s.aiChatStorage.SaveAIChatMessage(ctx, message)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to save message: %w", err)
	}

	//save planner answer
	err = s.aiChatStorage.SaveAIChatMessage(ctx, model.ChatMessage{
		Role:    model.RoleAssistant,
		Content: replyMessage,
		TripID:  message.TripID,
	})
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to save message: %w", err)
	}

	err = s.notifyUtils.FormAndSendNotifyMessage(
		ctx,
		message.TripID,
		"chat_reply_message",
		replyMessage,
		userID,
	)
	if err != nil {
		s.sendEventIfFailed(ctx, message, userID)
		return fmt.Errorf("failed to send chat_reply_message notify: %w", err)
	}

	return nil
}

func (s *AIChatService) sendEventIfFailed(ctx context.Context, message model.ChatMessage, userID int) {
	_ = s.notifyUtils.FormAndSendNotifyMessage(
		ctx,
		message.TripID,
		"chat_reply_message",
		"Что-то пошло не так, попробуйте снова",
		userID,
	)
}

func (s *AIChatService) formPromptOverUserMessage(
	userMessage model.ChatMessage,
	tripArea string,
) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Помоги спланировать поездку в город %s. ", tripArea))
	sb.WriteString(fmt.Sprintf("Вот вопрос пользователя: %s.\n", userMessage.Content))
	sb.WriteString("ФОРМАТ JSON ДОЛЖЕН СООТВЕТСТВОВАТЬ СЛЕДУЮЩЕЙ СТРУКТУРЕ: \n" +
		"{\n" +
		"\"places\" []Place\n" +
		"\"message\" string (тут должны быть твои слова)\n" +
		"}\n\n" +
		"type Place struct {\n" +
		"	\"name\" string\n" +
		"	\"recommended_visiting_time\" integer (кол-во часов)\n" +
		"}\n" +
		"НУЖНО ВЕРНУТЬ []Place (массив Place) БЕЗ ЛИШНИХ КОММЕНТАРИЕВ И БЕЗ ФОРМАТИРОВАНИЯ json.\n",
	//"ЕСЛИ ВОПРОС ПОЛЬЗОВАТЕЛЯ НЕ ОТНОСИТСЯ К ПЛАНИРОВАНИЮ ПОЕЗДКИ, А ТАКЖЕ МЕСТ, ГДЕ МОЖНО ПОКУШАТЬ, ОТВЕТЬ \"FAIL\"",
	)

	return sb.String(), nil
}

func (s *AIChatService) processPlannerResponse(ctx context.Context, plannerResponse, tripArea string) (string, error) {
	fmt.Println("resp: ", plannerResponse)
	//if plannerResponse == "FAIL" {
	//	return "Извините, кажется, данный вопрос не относится к планированию путешествия :)", nil
	//}

	type recommendedPlace struct {
		Name                    string `json:"name"`
		RecommendedVisitingTime int    `json:"recommended_visiting_time"`
	}
	type recommendedPlaceResponse struct {
		Places  []recommendedPlace `json:"places"`
		Message string             `json:"message"`
	}
	var parsedResponse recommendedPlaceResponse

	err := json.Unmarshal([]byte(plannerResponse), &parsedResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal planner response: %w", err)
	}

	var placesDomain []model.GooglePlace
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, recommendedPlace := range parsedResponse.Places {
		wg.Add(1)
		go func() {
			defer wg.Done()
			searchStr := fmt.Sprintf("%s %s", tripArea, recommendedPlace.Name)
			places, err := s.googleApi.FindPlace(ctx, searchStr, []string{
				"formatted_address",
				"name",
				"rating",
				"geometry",
				"photo",
				"place_id",
			})
			if err != nil {
				return
			}

			mu.Lock()
			placesDomain = append(placesDomain, places[0])
			mu.Unlock()
		}()
	}
	wg.Wait()

	type processPlannerResponse struct {
		Places  []model.GooglePlace `json:"places"`
		Message string              `json:"message"`
	}
	response := processPlannerResponse{
		Message: parsedResponse.Message,
		Places:  placesDomain,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(jsonData), nil
}
