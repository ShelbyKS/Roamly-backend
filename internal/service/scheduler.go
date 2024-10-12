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

func (s *SchedulerService) createPlaceList(places []*model.Place) string {
	var builder strings.Builder

	for i, place := range places {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i, place.Name))
	}

	return builder.String()
}

func (s *SchedulerService) GetSchedule(ctx context.Context, places []*model.Place) (model.Schedule, error) {
	listPlaces := s.createPlaceList(places)
	resp, err := s.client.PostPrompt(ctx, fmt.Sprintf(`СПЛАНИРУЙ ПОЕЗДКУ ТОЛЬКО ПО ДАННЫМ МЕСТАМ, ИСПОЛЬЗУЯ ИНФОРМАЦИЮ, ПРИВЕДЕННУЮ НИЖЕ, ВЕРНИ МЕСТА И ВРЕМЯ ПОСЕЩЕНИЯ, В ОТВЕТЕ ОПИШИ ИМЕННО ТОЛЬКО JSON  объект, КОТОРЫЙ БУДЕТ ОПИСЫВАТЬ СПЛАНИРОВАННОЕ РАСПИСАНИЕ, КРОМЕ JSON В ОТЕТЕ НИЧЕГО НЕ ДОЛЖНО БЫТЬ, МАРШУРТ ДОЛЖЕН БЫТЬ ОПТИМАЛЬНЫМ И УЧИТЫВАТЬ ВРЕМЯ НА ДОРОГУ МЕЖДУ МЕСТАМИ:
Дата поездки:
13.10-14.10
Места поездки:
%s
ФОРМАТ JSON ДОЛЖЕН СООТВЕТСТВОВАТЬ СЛЕДУЮЩЕЙ СТРУКТУРЕ: type Event struct {
  Place     string
  StartTime time.Time
  EndTime   time.Time
  Payload   map[string]any
}

type Schedule struct {
  Events  []Event
  Payload map[string]any
}`, listPlaces))
	if err != nil {
		return model.Schedule{}, err
	}

	var schedule model.Schedule

	err = json.Unmarshal([]byte(resp), &schedule)
	if err != nil {
		return model.Schedule{}, err
	}

	return schedule, err
}
