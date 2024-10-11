package service

import (
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
)

type UserService struct {
	userStorage storage.IUserStorage
}

func NewService(userStorage storage.IUserStorage) service.IUserService {
	return &UserService{
		userStorage: userStorage,
	}
}
