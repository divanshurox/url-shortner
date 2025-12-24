package config

import (
	"UrlShortener/internal/logger"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

type Config struct {
	// HTTP
	HttpPort string

	// Postgres
	DatabaseURL       string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxConnLifetime time.Duration
	DBMaxConnIdleTime time.Duration
	DBHealthTimeout   time.Duration

	// Redis
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	RedisHealthTimeout time.Duration
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getIntFromEnv(key string, def int32) int32 {
	val := getEnv(key, "")
	if val == "" {
		return def
	}
	var intVal int32
	_, err := fmt.Sscanf(val, "%d", intVal)
	if err != nil {
		logger.Logger().Fatal("invalid duration for", zap.String("key", key), zap.String("val", val))
	}
	return intVal
}

func getTimeFromEnv(key string, def time.Duration) time.Duration {
	val := getEnv(key, "")
	if val == "" {
		return def
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		logger.Logger().Fatal("invalid duration for", zap.String("key", key), zap.String("val", val))
	}
	return duration
}

func LoadConfig() *Config {
	return &Config{
		HttpPort: getEnv("HTTP_PORT", "8080"),

		DatabaseURL:       getEnv("DATABASE_URL", "postgres://app:secret@localhost:5432/urlshortener?sslmode=disable"),
		DBMaxConns:        getIntFromEnv("DB_MAX_CONNS", 20),
		DBMinConns:        getIntFromEnv("DB_MIN_CONNS", 2),
		DBMaxConnLifetime: getTimeFromEnv("DB_MAX_CONN_LIFETIME", time.Hour),
		DBMaxConnIdleTime: getTimeFromEnv("DB_MAX_CONN_IDLE", 30*time.Minute),
		DBHealthTimeout:   getTimeFromEnv("DB_HEALTH_TIMEOUT", 2*time.Second),

		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            0,
		RedisHealthTimeout: getTimeFromEnv("REDIS_HEALTH_TIMEOUT", 2*time.Second),
	}
}
