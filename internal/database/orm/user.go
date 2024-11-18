package orm

import "time"

type User struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Login     string
	Email     string
	Password  string
	Trips     []*Trip   `gorm:"many2many:trip_users;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
