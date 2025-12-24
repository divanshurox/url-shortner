package server

import (
	"UrlShortener/internal/cache"
	"UrlShortener/internal/config"
	"UrlShortener/internal/db"
	"UrlShortener/internal/handler"
	"UrlShortener/internal/logger"
	"UrlShortener/internal/routes"
	"UrlShortener/internal/service"
	"UrlShortener/internal/utils"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type GracefulServer struct {
	server   *http.Server
	listener net.Listener
	pg       *db.Postgres
	rd       *cache.Redis
}

type GracefulServerInterface interface {
	Prestart() error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func NewServer(ctx context.Context, cnf *config.Config) (*GracefulServer, error) {
	mux := chi.NewMux()
	err := logger.InitializeLogger()
	if err != nil {
		return nil, err
	}
	logger.Logger().Info("configuration", zap.Any("config", cnf))
	pg, err := db.InitPostgres(ctx, &db.PostgresOpts{
		DatabaseURL:     cnf.DatabaseURL,
		MaxConns:        cnf.DBMaxConns,
		MinConns:        cnf.DBMinConns,
		MaxConnLifetime: cnf.DBMaxConnLifetime,
		MaxConnIdleTime: cnf.DBMaxConnIdleTime,
		HealthTimeout:   cnf.DBHealthTimeout,
	})
	if err != nil {
		logger.Logger().Fatal("unable to initialize postgres client", zap.Error(err))
	}

	rd, err := cache.InitRedis(ctx, &cache.RedisOpts{
		Addr:          cnf.RedisAddr,
		Password:      cnf.RedisPassword,
		DB:            cnf.RedisDB,
		HealthTimeout: cnf.RedisHealthTimeout,
	})
	if err != nil {
		logger.Logger().Fatal("unable to initialize redis client", zap.Error(err))
	}

	urlshortenerService := service.NewUrlShortenerService(ctx, pg, rd)
	urlshortenerHandler := handler.NewUrlShortenerHandler(urlshortenerService)
	router := routes.NewUrlRouter(urlshortenerHandler)
	mux.Route("/url", router.GetUrlRouter)
	server := &http.Server{
		Addr:              ":" + cnf.HttpPort,
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}
	return &GracefulServer{server: server, pg: pg, rd: rd}, nil
}

func (gs *GracefulServer) Prestart() error {
	return nil
}

func (gs *GracefulServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", gs.server.Addr)
	if err != nil {
		return err
	}
	gs.listener = listener
	errCh := make(chan error, 1)
	utils.SafeGo(func() {
		logger.Logger().Info("Server is now listening!", zap.String("address", gs.server.Addr))
		if err := gs.server.Serve(gs.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	})

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = gs.Stop(shutdownCtx) // drain connections
		return <-errCh
	case err := <-errCh:
		return err
	}
}

func (gs *GracefulServer) Stop(ctx context.Context) error {
	logger.Close()
	shutdownContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	// closes postgres connection
	gs.pg.Close()

	// closes redis connection
	gs.rd.Close()

	// closes both server and listener
	if err := gs.server.Shutdown(shutdownContext); err != nil {
		gs.server.Close()
		return err
	}
	return nil
}
