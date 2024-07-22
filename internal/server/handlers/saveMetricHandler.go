package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// SaveMetricHandler save metric to storage
// Previous handler should add valid metric to context
func SaveMetricHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		metric := &apimodels.Metrics{}

		// use json and http on our own
		dec := json.NewDecoder(c.Request.Body)
		if err := dec.Decode(metric); err != nil {
			slog.Error("can't parse JSON", "error", err)
			http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// or just
		//if err := c.BindJSON(&metric); err != nil {
		//	slog.Error("can't parse JSON", "error", err)
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}

		switch metric.MType {
		case constants.GaugeStr:
			if metric.Value == nil {
				slog.Error("Empty metric value, float64 required")
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("empty metric value"))
				return
			}
			s.UpdateGauge(metric.ID, *metric.Value)
			slog.Info("value set", "new", s.GetGauge(metric.ID))
		case constants.CounterStr:
			if metric.Delta == nil {
				slog.Error("Empty metric delta, integer required")
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("empty metric delta"))
				return
			}
			delta := int(*metric.Delta)
			s.UpdateCounter(metric.ID, delta)
			slog.Info("value set", "new", s.GetCounter(metric.ID))
		default:
			slog.Error("Invalid type", "type", metric.MType)
			http.Error(c.Writer, fmt.Sprintf(
				"Invalid type, only %s and %s supported",
				constants.GaugeStr,
				constants.CounterStr,
			), http.StatusBadRequest)
			return
		}

		c.Set(contextMetricResponseKey, metric)
		c.Next()
	}
}
