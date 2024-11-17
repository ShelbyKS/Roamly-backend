package model

type UserTripRole int

const (
	Owner UserTripRole = iota
	Editor
	Reader
)
