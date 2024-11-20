package service

import (
	"context"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/google/uuid"
)

type IInviteService interface {
	GetTripInvitations(ctx context.Context, tripID uuid.UUID) ([]model.Invite, error)
	EnableInvitation(ctx context.Context, invite model.Invite) (model.Invite, error)
	DisableInvitation(ctx context.Context, invite model.Invite) error
	JoinTrip(ctx context.Context, inviteToken string, userID int) (uuid.UUID, error)
	UpdateMember(ctx context.Context, tripID uuid.UUID, userID int, access string) error
	DeleteMember(ctx context.Context, tripID uuid.UUID, userID int) error
}
