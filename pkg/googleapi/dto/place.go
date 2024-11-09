package dto

type GooglePlace struct {
	ID               string           `json:"id"`
	FormattedAddress string           `json:"formattedAddress"`
	DisplayName      DisplayName      `json:"displayName"`
	Rating           float64          `json:"rating"`
	Location         Location         `json:"location"`
	Photos           []Photo          `json:"photos"`
	EditorialSummary EditorialSummary `json:"editorialSummary"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Photo struct {
	Name string `json:"name"`
	// Height         int    `json:"height"`
	// Width          int    `json:"width"`
}

type DisplayName struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type EditorialSummary struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}
