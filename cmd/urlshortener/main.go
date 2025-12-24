package main

import (
	"UrlShortener/internal/config"
	"UrlShortener/internal/logger"
	"UrlShortener/internal/server"
	"context"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

func main() {
	cnf := config.LoadConfig()
	ctx := context.Background()
	shutdownCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv, err := server.NewServer(shutdownCtx, cnf)
	if err != nil {
		panic(err)
	}

	if err := srv.Start(shutdownCtx); err != nil {
		logger.Logger().Error("server exited", zap.Error(err))
	}
}
