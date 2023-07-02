package model

import "time"

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

type URL struct {
	UID       string    `json:"uid"`
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
