package service

import (
	"context"

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
	return service.userStorage.GetUserByID(ctx, id)
}

func (service *UserService) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	return service.userStorage.GetUserByLogin(ctx, login)
}

func (service *UserService) CreateUser(ctx context.Context, user model.User) error {
	return service.userStorage.CreateUser(ctx, user)
}

func (service *UserService) UpdateUser(ctx context.Context, user model.User) error {
	return service.userStorage.UpdateUser(ctx, user)
}
