package model

import "time"

type ShortenRequest struct {
	LongURL string `json:"long_url"`
	Title   string `json:"title"`
	Domain  string `json:"domain"`
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
	Title     string    `json:"title"`
	Domain    string    `json:"domain"`
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
