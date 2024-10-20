package orm

import "time"

type User struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Login     string
	Email     string
	Password  []byte
	Trips     []*Trip   `gorm:"many2many:trip_user;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
