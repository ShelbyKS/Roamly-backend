package model

import "time"

type Place struct {
	ID      string `json:"id"`
	Trips   []*Trip
	Name    string  `json:"name"`
	Photo   string  `json:"photo"`
	Rating  float32 `json:"rating"`
	Closing time.Time
	Opening time.Time
}
