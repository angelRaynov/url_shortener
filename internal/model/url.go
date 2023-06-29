package model

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	Link    string
	LongURL string
}

type ExpandRequest struct {
	ShortURL string `json:"short_url"`
}
