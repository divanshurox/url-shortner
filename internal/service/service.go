package service

import (
	"UrlShortener/internal/model"
	"context"
)

type UrlShortenerService interface {
	CreateShortUrl(ctx context.Context, req *model.CreateShortUrlReq) (*model.CreateShortUrlRes, error)
	GetUrlFromShortUrl(ctx context.Context, req *model.GetUrlFromShortUrlReq) (*model.GetUrlFromShortUrlRes, error)
}
