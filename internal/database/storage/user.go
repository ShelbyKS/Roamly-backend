package storage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/internal/domain"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"

	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) storage.IUserStorage {
	return &UserStorage{
		db: db,
	}
}

func (storage *UserStorage) GetUserByID(ctx context.Context, id int) (model.User, error) {
	user := orm.User{
		ID: id,
	}

	tx := storage.db.WithContext(ctx).First(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(domain.ErrUserNotFound, tx.Error)
	}

	if tx.Error != nil {
		return model.User{}, tx.Error
	}

	return model.User{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}, tx.Error
}

func (storage *UserStorage) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	user := orm.User{
		Login: login,
	}

	tx := storage.db.WithContext(ctx).First(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(domain.ErrUserNotFound, tx.Error)
	}

	if tx.Error != nil {
		return model.User{}, tx.Error
	}

	return model.User{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}, tx.Error
}

func (storage *UserStorage) CreateUser(ctx context.Context, user model.User) error {
	tx := storage.db.WithContext(ctx).Create(&orm.User{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	})

	return tx.Error
}

func (storage *UserStorage) UpdateUser(ctx context.Context, user model.User) error {
	tx := storage.db.WithContext(ctx).
		Model(&orm.User{ID: user.ID}).
		Updates(&orm.User{
			Login:    user.Login,
			Password: user.Password,
		})

	return tx.Error
}
