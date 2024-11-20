package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/service"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type InviteService struct {
	inviteStorage storage.IInviteStorage
	tripStorage   storage.ITripStorage
	jwtKey        string
}

func NewInviteService(inviteStorage storage.IInviteStorage, tripStorage storage.ITripStorage, jwtKey string) service.IInviteService {
	return &InviteService{
		inviteStorage: inviteStorage,
		tripStorage:   tripStorage,
		jwtKey:        jwtKey,
	}
}

func (s *InviteService) EnableInvitation(ctx context.Context, invite model.Invite) (model.Invite, error) {
	existingInvite, err1 := s.inviteStorage.GetInviteByTripAccess(ctx, invite)
	if err1 != nil && !errors.Is(err1, domain.ErrInviteNotFound) {
		return model.Invite{}, fmt.Errorf("failed to get existing invite: %w", err1)
	}

	fmt.Println("EXS INV: ", existingInvite)

	inviteToken, err := s.generateInviteToken(invite)
	if err != nil {
		return model.Invite{}, fmt.Errorf("failed to generate invite token: %w", err)
	}

	if errors.Is(err1, domain.ErrInviteNotFound) {
		invite.Enable = true
		invite.Token = inviteToken

		err = s.inviteStorage.CreateInvite(ctx, invite)
		if err != nil {
			return model.Invite{}, fmt.Errorf("failed to create invite: %w", err)
		}

		return invite, nil
	}

	existingInvite.Enable = true
	existingInvite.Token = inviteToken

	err = s.inviteStorage.UpdateInviteByTripAccess(ctx, existingInvite)
	if err != nil {
		return model.Invite{}, fmt.Errorf("failed to update existing invite: %w", err)
	}

	return existingInvite, nil
}

func (s *InviteService) generateInviteToken(invite model.Invite) (string, error) {
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"trip_id": invite.TripID.String(),
		"access":  invite.Access,
		"jti":     jti,
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtKey))
}

func (s *InviteService) GetTripInvitations(ctx context.Context, tripID uuid.UUID) ([]model.Invite, error) {
	invitations, err := s.inviteStorage.GetInvitesByTripID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invites from storage: %w", err)
	}

	return invitations, nil
}

func (s *InviteService) DisableInvitation(ctx context.Context, invite model.Invite) error {
	invite.Enable = false

	err := s.inviteStorage.UpdateInviteByTripAccess(ctx, invite)
	if err != nil {
		return fmt.Errorf("failed to update existing invite in storage: %w", err)
	}

	return nil
}

func (s *InviteService) JoinTrip(ctx context.Context, inviteToken string, userID int) (uuid.UUID, error) {
	invitation, err := s.inviteStorage.GetInviteByToken(ctx, inviteToken)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get invite from storage: %w", err)
	}

	for _, tripUsers := range invitation.Trip.Users {
		if userID == tripUsers.ID {
			return invitation.TripID, nil
		}
	}

	if !invitation.Enable {
		return uuid.Nil, domain.ErrInviteForbidden
	}

	err = s.inviteStorage.JoinTripByInvite(ctx, invitation, userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to join trip in storage: %w", err)
	}

	return invitation.TripID, nil
}

func (s *InviteService) UpdateMember(ctx context.Context, tripID uuid.UUID, userID int, access string) error {
	role, err := model.RoleFromString(access)
	if err != nil {
		return fmt.Errorf("invalid access role: %w", err)
	}

	//trip, err := s.tripStorage.GetTripByID(ctx, tripID)
	//if errors.Is(err, domain.ErrTripNotFound) {
	//	return err
	//}
	//if err != nil {
	//	return fmt.Errorf("fail to get trip from storage: %w", err)
	//}

	//todo: check that member is not owner

	//todo: need to forbid to update member to owner

	err = s.inviteStorage.UpdateMember(ctx, tripID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to update member role storage: %w", err)
	}

	return nil
}

func (s *InviteService) DeleteMember(ctx context.Context, tripID uuid.UUID, userID int) error {
	//trip, err := s.tripStorage.GetTripByID(ctx, tripID)
	//if errors.Is(err, domain.ErrTripNotFound) {
	//	return err
	//}
	//if err != nil {
	//	return fmt.Errorf("fail to get trip from storage: %w", err)
	//}

	//todo: check that member is not owner

	err := s.inviteStorage.DeleteMember(ctx, tripID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete member in storage: %w", err)
	}

	return nil
}
