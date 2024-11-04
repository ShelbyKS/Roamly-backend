package domain

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrTripNotFound      = errors.New("trip not found")
	ErrPlaceNotFound     = errors.New("place not found")
	ErrEventNotFound     = errors.New("event not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func GetStatusCodeByError(err error) int {
	switch err {
	case ErrUserNotFound, ErrTripNotFound, ErrPlaceNotFound, ErrEventNotFound:
		return http.StatusNotFound
	case ErrSessionNotFound, ErrWrongCredentials:
		return http.StatusUnauthorized
	case ErrUserAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
