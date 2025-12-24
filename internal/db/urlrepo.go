package db

import (
	"UrlShortener/internal/logger"
	"context"
	"go.uber.org/zap"
	"time"
)

type UrlRepo struct {
	Q Querier
}

func NewUrlRepo(q Querier) *UrlRepo {
	return &UrlRepo{q}
}

func (r *UrlRepo) Create(ctx context.Context, shortUrl, originalUrl string, expiresAt *time.Time) error {
	q := `
		INSERT INTO shorturlmappings (short_url,original_url,expires_at)
		VALUES ($1, $2, $3);
	`
	_, err := r.Q.Exec(ctx, q, shortUrl, originalUrl, expiresAt)
	if err != nil {
		logger.Logger().Error("error postgres client while creating entry", zap.Error(err))
		return err
	}
	return nil
}

func (r *UrlRepo) Get(ctx context.Context, shortUrl string) (*ShortUrlEntry, error) {
	q := `
		SELECT short_url, original_url, expires_at, created_at
		FROM shorturlmappings
		WHERE short_url = $1
	`
	var entry ShortUrlEntry
	err := r.Q.QueryRow(ctx, q, shortUrl).Scan(&entry.ShortUrl, &entry.OriginalUrl, &entry.ExpiresAt, &entry.CreatedAt)
	if err != nil {
		logger.Logger().Error("error postgres client while getting entry", zap.Error(err))
		return nil, err
	}
	return &entry, nil
}
