package model

import "time"

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int       `json:"-"`
}

type Credentials struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}
