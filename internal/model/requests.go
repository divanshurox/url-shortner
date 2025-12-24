package model

import "time"

type CreateShortUrlReq struct {
	Url            string    `json:"url"`
	Alias          string    `json:"alias,omitempty"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type CreateShortUrlRes struct {
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

type GetUrlFromShortUrlReq struct {
	ShortUrl string `json:"short_url"`
}

type GetUrlFromShortUrlRes struct {
	Url            string    `json:"url"`
	Alias          string    `json:"alias,omitempty"`
	ExpirationDate time.Time `json:"expiration_date"`
	CreationDate   time.Time `json:"creation_date"`
}
