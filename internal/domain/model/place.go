package model

type Place struct {
	ID     string  `json:"place_id"`
	Name   string  `json:"name"`
	Photo  string  `json:"photo"`
	Rating float32 `json:"rating"`
	Trips  []*Trip `json:"trips"`
}
