package model

type Trip struct {
	ID        int     `json:"id"`
	Users     []*User `json:"users"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Region    string  `json:"region"`
}
