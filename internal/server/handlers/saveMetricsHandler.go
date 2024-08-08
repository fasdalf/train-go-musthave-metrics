package handlers

import (
	"errors"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// SaveMetricsHandler save metric batch to a storage
func SaveMetricsHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		metrics := []apimodels.Metrics{}

		if err := c.BindJSON(&metrics); err != nil {
			slog.Error("can't parse JSON", "error", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("can't parse JSON"))
			return
		}

		err := s.SaveCommonModels(metrics)
		if err != nil {
			slog.Error("can't save metrics", "error", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("can't save metrics"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"savedCount": len(metrics)})
		c.Next()
	}
}
