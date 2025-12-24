package routes

import (
	"UrlShortener/internal/handler"
	"github.com/go-chi/chi/v5"
)

type UrlRouter struct {
	UrlHandler *handler.Handler
}

func NewUrlRouter(urlHandler *handler.Handler) *UrlRouter {
	return &UrlRouter{urlHandler}
}

func (ur *UrlRouter) GetUrlRouter(r chi.Router) {
	r.Post("/", ur.UrlHandler.CreateShortUrl)
	r.Get("/{short_url}", ur.UrlHandler.GetOriginalUrlFromShortUrl)
}
