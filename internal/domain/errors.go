package domain

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrTripNotFound     = errors.New("trip not found")
	ErrPlaceNotFound    = errors.New("place not found")
	ErrSessionNotFound  = errors.New("session not found")
	ErrWrongCredentials = errors.New("wrong credentials")
)
