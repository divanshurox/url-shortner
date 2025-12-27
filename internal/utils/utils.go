package utils

import (
	"UrlShortener/internal/logger"
	"go.uber.org/zap"
	"runtime/debug"
)

const (
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func SafeGo(fn func()) {
	go func() {
		if r := recover(); r != nil {
			logger.Logger().Error("panic in goroutine", zap.Any("panic", r), zap.ByteString("stack", debug.Stack()))
		}
		fn()
	}()
}

func EncodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}
	base := int64(62)
	output := make([]byte, 0, 11)
	for num > 0 {
		rem := num % base
		output = append(output, base62Chars[rem])
		num = num / base
	}
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}
	return string(output)
}
