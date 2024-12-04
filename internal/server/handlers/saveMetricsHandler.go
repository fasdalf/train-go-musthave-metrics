package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/catchable"
)

type Retryer interface {
	Try(ctx context.Context, do func() error, isRetryable func(err error) bool) (int, error)
}

// SaveMetricsHandler save metric batch to a storage
func SaveMetricsHandler(s Storage, retryer Retryer) func(c *gin.Context) {
	return func(c *gin.Context) {
		metrics := []apimodels.Metrics{}

		if err := c.BindJSON(&metrics); err != nil {
			slog.Error("can't parse JSON", "error", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("can't parse JSON"))
			return
		}

		if _, err := retryer.Try(c, func() error { return s.SaveCommonModels(c, metrics) }, catchable.IsPgConnectionError); err != nil {
			slog.Error("can't save metrics", "error", err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("can't save metrics"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"savedCount": len(metrics)})
		c.Next()
	}
}
