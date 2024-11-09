package dto

type TextSearchRequest struct {
	TextQuery    string       `json:"textQuery"`
	IncludeTypes string       `json:"includedType"`
	PageSize     int          `json:"pageSize"`
	LocationBias LocationBias `json:"locationBias"`
	LanguageCode string       `json:"languageCode"`
	PageToken    string       `json:"pageToken"`
}

type TextSearchResult struct {
	Results []GooglePlace `json:"places"`
}
