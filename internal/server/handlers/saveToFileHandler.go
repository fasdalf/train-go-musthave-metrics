package handlers

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/jsonofflinestorage"
	"github.com/gin-gonic/gin"
	"log/slog"
)

// NewSaveToFileHandler save all metrics to file
func NewSaveToFileHandler(ch jsonofflinestorage.SavedChan) func(c *gin.Context) {
	saveToFileHandler := func(c *gin.Context) {
		select {
		case ch <- struct{}{}:
			slog.Info("save signal sent")
		default:
			slog.Error("can't send save signal")
		}

		c.Next()
	}
	return saveToFileHandler
}
