package dto

type LocationRestriction struct {
	Circle Circle `json:"circle"`
}

type LocationBias struct {
	Circle Circle `json:"circle"`
}

type Circle struct {
	Center Center  `json:"center"`
	Radius float64 `json:"radius"`
}

type Center struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
