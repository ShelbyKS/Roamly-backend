package domain

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrTripNotFound       = errors.New("trip not found")
	ErrPlaceNotFound      = errors.New("place not found")
	ErrEventNotFound      = errors.New("event not found")
	ErrInviteNotFound     = errors.New("invite not found")
	ErrInviteForbidden    = errors.New("invite forbidden")
	ErrSessionNotFound    = errors.New("session not found")
	ErrWrongCredentials   = errors.New("wrong credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrPlaceAlreadyExists = errors.New("place already exists")
)

func GetStatusCodeByError(err error) int {
	switch err {
	case ErrUserNotFound, ErrTripNotFound, ErrPlaceNotFound, ErrEventNotFound, ErrInviteNotFound:
		return http.StatusNotFound
	case ErrInviteForbidden:
		return http.StatusForbidden
	case ErrSessionNotFound, ErrWrongCredentials:
		return http.StatusUnauthorized
	case ErrUserAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
