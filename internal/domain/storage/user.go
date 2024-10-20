package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IUserStorage interface {
	GetUserByID(ctx context.Context, id int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user model.User) error
}
