package service

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IAuthService interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	Login(ctx context.Context, user model.User) (model.Session, error)
	Logout(ctx context.Context, session model.Session) error
}
