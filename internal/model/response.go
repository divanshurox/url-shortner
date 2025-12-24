package model

import "time"

type CreateShortUrlResp struct {
	ShortUrl  string `json:"short_url"`
	CreatedAt string `json:"created_at"`
}

type GetUrlResp struct {
	Url            string    `json:"url"`
	ExpirationDate time.Time `json:"expiration_date"`
}
