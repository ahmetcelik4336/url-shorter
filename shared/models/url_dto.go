package dto

import "time"

type UrlResponse struct {
	Id                      int    `json:"id,omitempty"`
	ShortCode               string `json:"short_code,omitempty"`
	LongUrl                 string `json:"long_url,omitempty"`
	CreateAt                string `json:"created_at,omitempty"`
	Alias                   string `json:"alias,omitempty"`
	IsEncrypted             bool   `json:"is_encrypted,omitempty"`
	IsEncryptedText         string `json:"is_encrypted_text,omitempty"`
	ExpirationDate          string `json:"expiration_date"`
	ResultUrl               string `json:"url"`
	Password                string `json:"password"`
	ExpirationDateFormatted string `json:"dateformatted"`
}

type CreateUrlRequest struct {
	LongUrl        string    `json:"long_url" form:"url"`
	Alias          string    `json:"alias,omitempty" form:"alias"`
	Password       string    `json:"password,omitempty" form:"password"`
	ExpirationDate time.Time `json:"expiration_date,omitempty" form:"expiration"`
}

type UpdateUrlRequest struct {
	ID             int       // Bir önceki adımda parse ettiğiniz id
	LongUrl        string    `json:"long_url" form:"url"`
	Alias          string    `json:"alias,omitempty" form:"alias"`
	Password       string    `json:"password,omitempty" form:"password"`
	ExpirationDate time.Time `json:"expiration_date,omitempty" form:"expiration"`
}

type UrlConfig struct {
	IsRedirect bool
}

type ExcelUrlRequest struct {
	LongUrl        string     `json:"long_url" form:"url"`
	Alias          string     `json:"alias,omitempty" form:"alias"`
	Password       *string    `json:"password,omitempty" form:"password"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty" form:"expiration"`
}
