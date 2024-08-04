package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

// NewSaveToFileHandler save all metrics to file
func NewSaveToFileHandler(s FileStorage) func(c *gin.Context) {
	saveToFileHandler := func(c *gin.Context) {
		err := s.SaveWithInterval()
		if err != nil {
			slog.Error("can't save to file", "error", err)
		}

		c.Next()
	}
	return saveToFileHandler
}
