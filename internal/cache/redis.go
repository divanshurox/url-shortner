package cache

import (
	"UrlShortener/internal/logger"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type Redis struct {
	Client *redis.Client
}

type RedisOpts struct {
	Addr          string
	Password      string
	DB            int
	HealthTimeout time.Duration
}

func InitRedis(ctx context.Context, opts *RedisOpts) (*Redis, error) {
	r := redis.NewClient(&redis.Options{
		Password: opts.Password,
		Addr:     opts.Addr,
		DB:       opts.DB,
	})

	cctx, cancel := context.WithTimeout(ctx, opts.HealthTimeout)
	defer cancel()
	if err := r.Ping(cctx).Err(); err != nil {
		logger.Logger().Error("closing redis client", zap.Error(err))
		_ = r.Close()
		return nil, err
	}
	return &Redis{r}, nil
}

func (r *Redis) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}
