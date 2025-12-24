package db

import (
	"UrlShortener/internal/logger"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

type PostgresOpts struct {
	DatabaseURL     string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	HealthTimeout   time.Duration
}

func InitPostgres(ctx context.Context, opts *PostgresOpts) (*Postgres, error) {
	cfg, err := pgxpool.ParseConfig(opts.DatabaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = opts.MaxConns
	cfg.MinConns = opts.MinConns
	cfg.MaxConnIdleTime = opts.MaxConnIdleTime
	cfg.MaxConnLifetime = opts.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	cctx, cancel := context.WithTimeout(ctx, opts.HealthTimeout)
	defer cancel()
	if err = pool.Ping(cctx); err != nil {
		logger.Logger().Error("closing postgres client", zap.Error(err))
		pool.Close()
		return nil, err
	}

	return &Postgres{pool}, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
