package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/gin-gonic/gin"
	"html"
	"log/slog"
	"net/http"
)

func NewViewStatsHandler(ms metricstorage.Storage) func(c *gin.Context) {
	// In our project we have a tradition to add single middleware to put metricstorage.Storage pointer in context.
	// Let's use handler constructors for now
	return func(c *gin.Context) {
		mType := c.Param(`type`)
		mName := c.Param(`name`)
		mValue := ``
		switch mType {
		case constants.GaugeStr:
			if !ms.HasGauge(mName) {
				_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf(`metric "%s" not found`, html.EscapeString(mName)))
			}
			mValue = fmt.Sprint(ms.GetGauge(mName))
		case constants.CounterStr:
			if !ms.HasCounter(mName) {
				_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf(`metric "%s" not found`, html.EscapeString(mName)))
			}
			mValue = fmt.Sprint(ms.GetCounter(mName))
		default:
			slog.Error("Invalid type", "type", mType)
			http.Error(c.Writer, fmt.Sprintf(
				"Invalid type, only %s and %s supported",
				constants.GaugeStr,
				constants.CounterStr,
			), http.StatusBadRequest)
			return
		}

		c.Header(`Content-Type`, `text/plain`)
		_, _ = c.Writer.Write([]byte(mValue))
		slog.Info("got value", "type", mType, "name", mName, "value", mValue)
	}
}
