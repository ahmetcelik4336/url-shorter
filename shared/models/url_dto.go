package dto

import "time"

type UrlResponse struct {
	Id        int    `json:"id,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
	LongUrl   string `json:"long_url,omitempty"`
	CreateAt  string `json:"created_at,omitempty"`
}

type CreateUrlRequest struct {
	LongUrl        string    `json:"long_url"`
	Alias          string    `json:"alias"`
	Password       string    `json:"password"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type UpdateUrlRequest struct {
	LongUrl string `json:"long_url"`
	ID      int    `json:"id"`
}

type UrlConfig struct {
	IsRedirect bool
}
