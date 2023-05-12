package model

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	Link    string
	LongURL string
}
