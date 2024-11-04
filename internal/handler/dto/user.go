package dto

import "time"

type GetUser struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
