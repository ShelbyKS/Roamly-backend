package service

import (
	"context"
	"fmt"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type UserService struct {
	userStorage storage.IUserStorage
}

func NewUserService(userStorage storage.IUserStorage) service.IUserService {
	return &UserService{
		userStorage: userStorage,
	}
}

func (service *UserService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	user, err := service.userStorage.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("fail to get user from storage: %w", err)
	}

	return user, nil
}

func (service *UserService) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	user, err := service.userStorage.GetUserByLogin(ctx, login)
	if err != nil {
		return model.User{}, fmt.Errorf("fail to get user from storage: %w", err)
	}

	return user, nil
}

func (service *UserService) CreateUser(ctx context.Context, user model.User) error {
	err := service.userStorage.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("fail to create user from storage: %w", err)
	}

	return nil
}

func (service *UserService) UpdateUser(ctx context.Context, user model.User) error {
	err := service.userStorage.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("fail to update user from storage: %w", err)
	}

	return nil
}
