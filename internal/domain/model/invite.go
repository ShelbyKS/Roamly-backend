package model

import "github.com/google/uuid"

type Invite struct {
	Token  string
	TripID uuid.UUID
	Trip   Trip
	Access string
	Enable bool
}
