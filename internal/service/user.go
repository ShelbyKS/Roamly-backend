package service

import (
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
