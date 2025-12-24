package utils

import (
	"UrlShortener/internal/logger"
	"go.uber.org/zap"
	"runtime/debug"
)

func SafeGo(fn func()) {
	go func() {
		if r := recover(); r != nil {
			logger.Logger().Error("panic in goroutine", zap.Any("panic", r), zap.ByteString("stack", debug.Stack()))
		}
		fn()
	}()
}
