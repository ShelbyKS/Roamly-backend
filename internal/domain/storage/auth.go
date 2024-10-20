package storage

import (
	"context"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type ISessionStorage interface {
	Add(ctx context.Context, session model.Session) error
	DeleteByToken(ctx context.Context, token string) error
	SessionExists(ctx context.Context, token string) (model.Session, error)
}
