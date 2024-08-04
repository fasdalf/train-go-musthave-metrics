package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

// SlogHandler gin middleware
func SlogHandler(c *gin.Context) {
	start := time.Now()

	responseData := &loggingResponseData{
		status: 0,
		size:   0,
	}
	lw := &loggingResponseWriter{
		ResponseWriter: c.Writer, // встраиваем оригинальный gin.ResponseWriter
		responseData:   responseData,
	}

	c.Writer = lw
	// передаём управление хендлеру
	c.Next()

	duration := time.Since(start)

	slog.Info("HTTP request processed",
		"uri", c.Request.RequestURI,
		"method", c.Request.Method,
		"status", responseData.status, // получаем перехваченный код статуса ответа
		"duration", duration.String(),
		"size", responseData.size, // получаем перехваченный размер ответа
	)
}
