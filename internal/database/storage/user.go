package storage

import (
	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/internal/domain/storage"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) storage.IUserStorage {
	return &UserStorage{
		db: db,
	}
}
