package handlers

import (
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"html"
	"log/slog"
	"net/http"
)

// NewViewStatsHandler view single stored metric
//
// #Deprecated: for old lesson
func NewViewStatsHandler(ms Storage) func(c *gin.Context) {
	// In our project we have a tradition to add single middleware to put metricstorage.Storage pointer in context.
	// Let's use handler constructors for now
	return func(c *gin.Context) {
		mType := c.Param(URLParamType)
		mName := c.Param(URLParamName)
		mValue := ""
		switch mType {
		case constants.GaugeStr:
			if h, err := ms.HasGauge(mName); err != nil || !h {
				if err != nil {
					slog.Error("can't get gauge", "id", mName, "error", err)
					http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
					return
				}
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(mName)), http.StatusNotFound)
				return
			}
			value, err := ms.GetGauge(mName)
			if err != nil {
				slog.Error("can't get gauge", "name", mName, "error", err)
				_ = c.AbortWithError(http.StatusInternalServerError, errors.New(`unexpected error`))
				return
			}
			mValue = fmt.Sprint(value)
		case constants.CounterStr:
			if h, err := ms.HasCounter(mName); err != nil || !h {
				if err != nil {
					slog.Error("can't get gauge", "id", mName, "error", err)
					http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
					return
				}
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(mName)), http.StatusNotFound)
				return
			}
			value, err := ms.GetCounter(mName)
			if err != nil {
				slog.Error("can't get counter", "name", mName, "error", err)
				_ = c.AbortWithError(http.StatusInternalServerError, errors.New(`unexpected error`))
				return
			}
			mValue = fmt.Sprint(value)
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
