package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IInviteStorage interface {
	GetInviteByTripAccess(ctx context.Context, invite model.Invite) (model.Invite, error)
	CreateInvite(ctx context.Context, invite model.Invite) error
	UpdateInviteByTripAccess(ctx context.Context, invite model.Invite) error
	GetInvitesByTripID(ctx context.Context, tripID uuid.UUID) ([]model.Invite, error)
	GetInviteByToken(ctx context.Context, token string) (model.Invite, error)
	JoinTripByInvite(ctx context.Context, invite model.Invite, userID int) error
}
