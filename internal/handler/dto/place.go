package dto

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Photo struct {
	PhotoReference string `json:"photo_reference"`
}

type GooglePlace struct {
	FormattedAddress    string   `json:"formatted_address"`
	Geometry            Geometry `json:"geometry"`
	Name                string   `json:"name"`
	Photos              []Photo  `json:"photos"`
	PlaceID             string   `json:"place_id"`
	Rating              float64  `json:"rating"`
	Types               []string `json:"types"`
	Vicinity            string   `json:"vicinity"`
	EditorialSummary    string   `json:"editorial_summary"`
	RecommendedDuration int      `json:"recommended_duration"`
}
