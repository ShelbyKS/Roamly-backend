package dto

type PlacesNearbyRequest struct {
	IncludeTypes        []string            `json:"includedTypes"`
	MaxResultCount      int                 `json:"maxResultCount"`
	RankPreference      string              `json:"rankPreference"`
	LocationRestriction LocationRestriction `json:"locationRestriction"`
	LanguageCode        string              `json:"languageCode"`
}

type LocationRestriction struct {
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

type PlacesNearbyResult struct {
	Results []GooglePlace `json:"places"`
}

type DisplayName struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type GooglePlace struct {
	FormattedAddress string      `json:"formattedAddress"`
	DisplayName      DisplayName `json:"displayName"`
	Rating           float64     `json:"rating"`
	Location         Location    `json:"location"`
	Photos           []Photo     `json:"photos"`
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
