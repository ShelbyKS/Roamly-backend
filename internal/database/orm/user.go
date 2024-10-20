package orm

import "time"

type User struct {
	ID        int `gorm:"primaryKey"`
	Login     string
	Email     string
	Password  []byte
	Trips     []*Trip   `gorm:"many2many:trip_users"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
