package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type SchedulerService struct {
	client storage.ISchedulerClient
}

func NewShedulerService(client storage.ISchedulerClient) service.ISchedulerService {
	return &SchedulerService{
		client: client,
	}
}

func (s *SchedulerService) GetSchedule(ctx context.Context, trip model.Trip, places []*model.Place, timeMatrix [][]int) (model.Schedule, error) {
	prompt, err := s.generateRequestString(trip, places, timeMatrix)
	if err != nil {
		return model.Schedule{}, err
	}

	fmt.Println(prompt)

	resp, err := s.client.PostPrompt(ctx, prompt)
	if err != nil {
		return model.Schedule{}, err
	}

	return ParseResponse(resp)
}

func (s *SchedulerService) generateRequestString(trip model.Trip, places []*model.Place, timeMatrix [][]int) (string, error) {
	var sb strings.Builder

	sb.WriteString(`СПЛАНИРУЙ ПОЕЗДКУ ТОЛЬКО ПО ДАННЫМ МЕСТАМ,
ИСПОЛЬЗУЯ ИНФОРМАЦИЮ, ПРИВЕДЕННУЮ НИЖЕ, ВЕРНИ МЕСТА И ВРЕМЯ ПОСЕЩЕНИЯ,
В ОТВЕТЕ ОПИШИ ИМЕННО ТОЛЬКО JSON  объект, КОТОРЫЙ БУДЕТ ОПИСЫВАТЬ СПЛАНИРОВАННОЕ РАСПИСАНИЕ,
КРОМЕ JSON В ОТЕТЕ НИЧЕГО НЕ ДОЛЖНО БЫТЬ, МАРШУРТ ДОЛЖЕН БЫТЬ ОПТИМАЛЬНЫМ И УЧИТЫВАТЬ ВРЕМЯ НА ДОРОГУ МЕЖДУ МЕСТАМИ,
КОТОРЫЕ УКАЗАНЫ В МАТРИЦЕ ВРЕМЕНИ, ГДЕ IJ-ОМУ СТОЛБЦУ СООТВЕТСТВУЕТ ВРЕМЯ ПУТИ ИЗ I В J:`)

	sb.WriteString("Дата поездки:\n")
	sb.WriteString(fmt.Sprintf("С: %s\n", trip.StartTime))
	sb.WriteString(fmt.Sprintf("По: %s\n", trip.EndTime))

	sb.WriteString("Места поездки:\n")
	for i, place := range places {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, place.GooglePlace.Name))
	}

	//sb.WriteString("Время работы:\n")
	//for i, place := range places {
	//	sb.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, place.Opening, place.Closing))
	//}

	sb.WriteString("Матрица времен:\n")
	for i := range timeMatrix {
		for j := range timeMatrix[i] {
			sb.WriteString(fmt.Sprintf("%d ", timeMatrix[i][j]))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(`
	ФОРМАТ JSON ДОЛЖЕН СООТВЕТСТВОВАТЬ СЛЕДУЮЩЕЙ СТРУКТУРЕ: type Event struct {
	PlaceID string    gorm:"primaryKey"
	TripID  uuid.UUID gorm:"primaryKey"
	Place   Place
	Trip    Trip

	StartTime string gorm:"type:TIME"
	EndTime   string gorm:"type:TIME"
	Payload   datatypes.JSON
}

	type Schedule struct {
		Events  []Event
		Payload map[string]any
	}

В ОТВЕТЕ ВЕРНИ ТОЛЬКО JSON  объект, КОТОРЫЙ БУДЕТ ОПИСЫВАТЬ СПЛАНИРОВАННОЕ РАСПИСАНИЕ. БЕЗ ЛИШНИХ КОММЕНТАРИЕВ.
`)
	return sb.String(), nil
}

func ParseResponse(response string) (model.Schedule, error) {
	var schedule model.Schedule

	err := json.Unmarshal([]byte(response), &schedule)
	if err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}
