package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
)

const (
	contextMetricResponseKey = "apimodels.Metrics"
)

// MetricValueResponseHandler respond with current metric value
// Previous handler should add valid metric to context
func MetricValueResponseHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var metric *apimodels.Metrics
		metricWrapped, ok := c.Get(contextMetricResponseKey)
		if !ok {
			slog.Error("No value in context", "key", contextMetricResponseKey)
			c.Next()
			return
		}

		metric, ok = metricWrapped.(*apimodels.Metrics)
		if !ok {
			slog.Error("value in context is not a metric", "key", contextMetricResponseKey, "type", reflect.TypeOf(metricWrapped))
			c.Next()
			return
		}

		switch metric.MType {
		case constants.GaugeStr:
			gauge, err := s.GetGauge(metric.ID)
			if err != nil {
				slog.Error("can't get gauge", "id", metric.ID, "error", err)
				http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
				return
			}
			metric.Value = &gauge
		case constants.CounterStr:
			counter, err := s.GetCounter(metric.ID)
			if err != nil {
				slog.Error("can't get counter", "id", metric.ID, "error", err)
				http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
				return
			}
			counter64 := int64(counter)
			metric.Delta = &counter64
		default:
			c.Next()
			return
		}

		// Have to set headers *before* writing to body
		c.Header("Content-Type", "application/json")
		// IRL just use c.IndentedJSON(200, metric)
		// Use encoder manually
		enc := json.NewEncoder(c.Writer)
		enc.SetIndent("", "  ")
		if err := enc.Encode(metric); err != nil {
			slog.Error(err.Error())
		}

		c.Next()
		slog.Info("Processed OK")
	}
}
