package dto

type PlacesNearbyRequest struct {
	IncludeTypes        []string            `json:"includedTypes"`
	MaxResultCount      int                 `json:"maxResultCount"`
	RankPreference      string              `json:"rankPreference"`
	LocationRestriction LocationRestriction `json:"locationRestriction"`
	LanguageCode        string              `json:"languageCode"`
}

type PlacesNearbyResult struct {
	Results []GooglePlace `json:"places"`
}
