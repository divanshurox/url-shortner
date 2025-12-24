package service

import (
	"UrlShortener/internal/cache"
	"UrlShortener/internal/db"
	"UrlShortener/internal/logger"
	"UrlShortener/internal/model"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type UrlShortenerServiceImpl struct {
	UrlRepo  *db.UrlRepo
	UrlCache *cache.ShortUrlCache
}

func NewUrlShortenerService(ctx context.Context, pg *db.Postgres, rd *cache.Redis) *UrlShortenerServiceImpl {
	urlRepo := db.NewUrlRepo(pg.Pool)
	stats := pg.Pool.Stat()
	logger.Logger().Info("pg pool",
		zap.Int32("totalConns", stats.TotalConns()),
		zap.Int32("idleConns", stats.IdleConns()),
	)
	urlCache := cache.NewShortUrlCache(rd)
	return &UrlShortenerServiceImpl{
		UrlRepo:  urlRepo,
		UrlCache: urlCache,
	}
}

func (s *UrlShortenerServiceImpl) CreateShortUrl(ctx context.Context, req *model.CreateShortUrlReq) (*model.CreateShortUrlRes, error) {
	shortUrl := generateShortUrl()
	if err := s.UrlRepo.Create(ctx, shortUrl, req.Url, &req.ExpirationDate); err != nil {
		return nil, err
	}
	res := &model.CreateShortUrlRes{
		ShortUrl:     shortUrl,
		CreationDate: time.Now(),
	}
	return res, nil
}

func (s *UrlShortenerServiceImpl) GetUrlFromShortUrl(ctx context.Context, req *model.GetUrlFromShortUrlReq) (*model.GetUrlFromShortUrlRes, error) {
	shortUrl := req.ShortUrl
	url, err := s.UrlCache.Get(ctx, shortUrl)
	if err != nil || url == "" {
		logger.Logger().Info("unable to get shorturl from cache", zap.String("shortUrl", shortUrl), zap.Error(err))
		row, err := s.UrlRepo.Get(ctx, shortUrl)
		if err != nil {
			logger.Logger().Error("error while getting shorturl entry from DB", zap.String("shortUrl", shortUrl), zap.Error(err))
			return nil, err
		}
		err = s.UrlCache.Set(ctx, shortUrl, row.OriginalUrl, 30*time.Minute)
		if err != nil {
			logger.Logger().Error("unable to set shorturl from cache", zap.String("shortUrl", shortUrl), zap.Error(err))
		}
		res := &model.GetUrlFromShortUrlRes{
			Url:            row.OriginalUrl,
			Alias:          "",
			ExpirationDate: row.ExpiresAt,
			CreationDate:   row.CreatedAt,
		}
		return res, nil
	}
	logger.Logger().Info("got shorturl from cache", zap.String("shortUrl", shortUrl), zap.String("original_url", url))
	res := &model.GetUrlFromShortUrlRes{
		Url:            url,
		Alias:          "",
		ExpirationDate: time.Time{},
		CreationDate:   time.Time{},
	}
	return res, nil
}

func generateShortUrl() string {
	return uuid.New().String()
}
