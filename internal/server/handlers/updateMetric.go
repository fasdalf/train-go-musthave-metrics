package handlers

import (
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"html"
	"log/slog"
	"net/http"
	"strconv"
)

// NewUpdateMetricHandler update single metric value
func NewUpdateMetricHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		mType := c.Param(`type`)
		mName := c.Param(`name`)
		mValue := c.Param(`value`)
		slog.Info("Processing update", "type", mType, "name", mName, "value", mValue)

		if mName == "" {
			slog.Error("Metric not found")
			_ = c.AbortWithError(http.StatusNotFound, errors.New("Metric not found"))
			return
		}

		switch mType {
		case constants.GaugeStr:
			floatValue, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				slog.Error("Invalid metric value, float64 required", `value`, html.EscapeString(mType))
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid metric value"))
				return
			}
			s.UpdateGauge(mName, floatValue)
			slog.Info("value set", "new", s.GetGauge(mName))
		case constants.CounterStr:
			intValue, err := strconv.Atoi(mValue)
			if err != nil {
				slog.Error("Invalid metric value, integer required", `value`, html.EscapeString(mType))
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid metric value"))
				return
			}
			s.UpdateCounter(mName, intValue)
			slog.Info("value set", "new", s.GetGauge(mName))
		default:
			slog.Error("Type is invalid", "type", html.EscapeString(mType))
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf(
				`Invalid type "%s", only "%s" and "%s" are supported`,
				html.EscapeString(mType),
				constants.GaugeStr,
				constants.CounterStr,
			))
			return
		}

		c.Next()
		slog.Info("Processed OK")
	}
}
