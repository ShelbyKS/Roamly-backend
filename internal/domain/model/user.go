package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
