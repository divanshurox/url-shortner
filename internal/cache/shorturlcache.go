package cache

import (
	"context"
	"time"
)

type ShortUrlCache struct {
	*Redis
}

func NewShortUrlCache(redis *Redis) *ShortUrlCache {
	return &ShortUrlCache{redis}
}

func (s *ShortUrlCache) Get(ctx context.Context, shortUrl string) (string, error) {
	return s.Client.Get(ctx, shortUrl).Result()
}

func (s *ShortUrlCache) Set(ctx context.Context, shortUrl, longUrl string, ttl time.Duration) error {
	return s.Client.Set(ctx, shortUrl, longUrl, ttl).Err()
}
