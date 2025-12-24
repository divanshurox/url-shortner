package db

import "time"

type ShortUrlEntry struct {
	ShortUrl    string
	OriginalUrl string
	ExpiresAt   time.Time
	CreatedAt   time.Time
}
