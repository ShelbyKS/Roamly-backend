package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IUserService interface {
	GetUserByID(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
}
