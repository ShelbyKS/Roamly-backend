package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/google/uuid"

	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type InviteStorage struct {
	db *gorm.DB
}

func NewInviteStorage(db *gorm.DB) storage.IInviteStorage {
	return &InviteStorage{
		db: db,
	}
}

func (storage *InviteStorage) GetInviteByTripAccess(ctx context.Context, invite model.Invite) (model.Invite, error) {
	inviteDB := &orm.Invite{}

	tx := storage.db.WithContext(ctx).
		Model(&orm.Invite{}).
		Where("trip_id = ? AND access = ?", invite.TripID, invite.Access).
		First(inviteDB)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.Invite{}, domain.ErrInviteNotFound
	}

	if tx.Error != nil {
		return model.Invite{}, tx.Error
	}

	return InviteConverter{}.ToDomain(*inviteDB), nil
}

func (storage *InviteStorage) CreateInvite(ctx context.Context, invite model.Invite) error {
	inviteDb := InviteConverter{}.ToDb(invite)

	tx := storage.db.WithContext(ctx).Create(&inviteDb)

	return tx.Error
}

func (storage *InviteStorage) UpdateInviteByTripAccess(ctx context.Context, invite model.Invite) error {
	inviteDb := InviteConverter{}.ToDb(invite)

	tx := storage.db.WithContext(ctx).
		Model(&orm.Invite{}).
		Where("trip_id = ? AND access = ?", invite.TripID, invite.Access).
		Updates(&inviteDb)

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return domain.ErrInviteNotFound
	}

	return nil
}

func (storage *InviteStorage) GetInvitesByTripID(ctx context.Context, tripID uuid.UUID) ([]model.Invite, error) {
	var invitesDB []orm.Invite

	tx := storage.db.WithContext(ctx).
		Where("trip_id = ? AND enable = true", tripID).
		Find(&invitesDB)

	if tx.Error != nil {
		return []model.Invite{}, tx.Error
	}

	var invites []model.Invite
	for _, invite := range invitesDB {
		invites = append(invites, InviteConverter{}.ToDomain(invite))
	}

	return invites, nil
}

func (storage *InviteStorage) GetInviteByToken(ctx context.Context, token string) (model.Invite, error) {
	inviteDB := &orm.Invite{}

	fmt.Println("TOKEN:", token)

	tx := storage.db.WithContext(ctx).
		Model(&orm.Invite{}).
		Where("token = ?", token).
		Preload("Trip").
		Preload("Trip.Users").
		First(inviteDB)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.Invite{}, domain.ErrInviteNotFound
	}

	if tx.Error != nil {
		return model.Invite{}, tx.Error
	}

	return InviteConverter{}.ToDomain(*inviteDB), nil
}

func (storage *InviteStorage) JoinTripByInvite(ctx context.Context, invite model.Invite, userID int) error {
	userRole, err := model.RoleFromString(invite.Access)
	if err != nil {
		return err
	}

	tripUser := orm.TripUsers{
		UserID:   userID,
		TripID:   invite.TripID,
		UserRole: int(userRole),
	}

	tx := storage.db.WithContext(ctx).Create(&tripUser)

	return tx.Error
}
