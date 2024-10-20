package postgresql

import (
	"context"
	"errors"
	"fmt"

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

func (storage *UserStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user := orm.User{}

	res := storage.db.WithContext(ctx).
		Where(&orm.User{Email: email}).
		First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, domain.ErrUserNotFound
	}

	if res.Error != nil {
		return model.User{}, fmt.Errorf("failed to get user by email: %w", res.Error)
	}

	return model.User{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}, nil
}

func (storage *UserStorage) CreateUser(ctx context.Context, user *model.User) error {
	usrModel := orm.User{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
	}

	res := storage.db.WithContext(ctx).Create(&usrModel)
	if res.Error != nil {
		return fmt.Errorf("failed to create user: %s", res.Error)
	}

	user.ID = usrModel.ID
	user.CreatedAt = usrModel.CreatedAt

	return res.Error
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
