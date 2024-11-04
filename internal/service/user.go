package service

import (
	"context"
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type UserService struct {
	userStorage    storage.IUserStorage
	sessionStorage storage.ISessionStorage
}

func NewUserService(userStorage storage.IUserStorage, sessionStorage storage.ISessionStorage) service.IUserService {
	return &UserService{
		userStorage:    userStorage,
		sessionStorage: sessionStorage,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	user, err := s.userStorage.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("fail to get user from storage: %w", err)
	}

	return user, nil
}

func (service *UserService) UpdateUser(ctx context.Context, user model.User) error {
	err := service.userStorage.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("fail to update user from storage: %w", err)
	}

	return nil
}
