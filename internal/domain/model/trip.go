package model

type Trip struct {
	ID        int     `json:"id"`
	Users     []*User `json:"users"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	AreaID    string  `json:"area_id"`
	Area      *Place
	Places    []*Place `json:"places"`
}
