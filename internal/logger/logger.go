package logger

import (
	"go.uber.org/zap"
	"sync"
)

var (
	loggerInstance *zap.Logger
	loggerError    error
	once           sync.Once
)

func InitializeLogger() error {
	once.Do(func() {
		loggerInstance, loggerError = zap.NewDevelopment()
	})
	return loggerError
}

func Logger() *zap.Logger {
	return loggerInstance
}

func Close() {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Sync()
	loggerInstance = nil
}
