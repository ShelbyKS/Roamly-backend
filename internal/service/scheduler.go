package service

import (
	"context"
	"time"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type SchedulerService struct {
	client storage.ISchedulerClient
}

func (s *SchedulerService) GetSchedule(ctx context.Context, places []model.Place) (model.Schedule, error) {
	resp, err := s.client.PostPrompt(ctx, "Хэллоу ворлд")
	if err != nil {
		return model.Schedule{}, err
	}

	return model.Schedule{
		Events: []model.Event{
			model.Event{
				Place: "Москва",
				StartTime: time.Now(),
				EndTime: time.Now(),
				Payload: map[string]any{
					"Kek": "Lol",
				},
			},
		},
		Payload: map[string]any{
			"response": resp,
		},
	}, nil
}
