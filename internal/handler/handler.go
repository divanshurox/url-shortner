package handler

import (
	"UrlShortener/internal/logger"
	"UrlShortener/internal/model"
	"UrlShortener/internal/service"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	service service.UrlShortenerService
}

func NewUrlShortenerHandler(service service.UrlShortenerService) *Handler {
	return &Handler{
		service,
	}
}

// POST: /url
func (h *Handler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.CreateShortUrlReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger().Error("failed to decode request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	logger.Logger().Info("CreateShortUrl", zap.Any("req", req))
	url := req.Url
	expirationDate := req.ExpirationDate
	if url == "" || expirationDate.IsZero() {
		logger.Logger().Error("failed to decode request body", zap.String("url", url), zap.Any("expirationDate", expirationDate))
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}
	if !validUrl(url) {
		logger.Logger().Error("url not valid", zap.String("url", url))
		http.Error(w, "url not valid", http.StatusBadRequest)
		return
	}
	res, err := h.service.CreateShortUrl(ctx, &req)
	if err != nil {
		logger.Logger().Error(err.Error(), zap.String("url", url), zap.Any("expirationDate", expirationDate))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(res); err != nil {
		logger.Logger().Error("failed to encode response body", zap.Error(err))
		http.Error(w, "invalid response body", http.StatusInternalServerError)
		return
	}
}

func validUrl(url string) bool {
	return true
}

// GET /url/{short_url}
func (h *Handler) GetOriginalUrlFromShortUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortUrl := r.PathValue("short_url")
	if shortUrl == "" {
		http.Error(w, "missing short url code", http.StatusBadRequest)
		return
	}
	req := &model.GetUrlFromShortUrlReq{ShortUrl: shortUrl}
	res, err := h.service.GetUrlFromShortUrl(ctx, req)
	if err != nil {
		logger.Logger().Error(err.Error(), zap.String("short_url", req.ShortUrl))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, res.Url, http.StatusFound)
	return
}
