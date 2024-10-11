package user

import (
	"gorm.io/gorm"
)

type Storage interface {
	//getByID(id uint64) (domain.User, error)
	//create(user domain.User) (bool, error)
	//update(user domain.User) (bool, error)
}

type storageImpl struct {
	conn *gorm.DB
}

func NewStorage(conn *gorm.DB) Storage {
	return &storageImpl{
		conn: conn,
	}
}
