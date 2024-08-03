package handlers

import (
	"encoding/json"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"log/slog"
	"reflect"
)

const (
	contextMetricResponseKey = "apimodels.Metrics"
)

// MetricValueResponseHandler respond with current metric value
// Previous handler should add valid metric to context
func MetricValueResponseHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var metric *apimodels.Metrics = nil
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
			gauge := s.GetGauge(metric.ID)
			metric.Value = &gauge
		case constants.CounterStr:
			counter := int64(s.GetCounter(metric.ID))
			metric.Delta = &counter
		default:
			c.Next()
			return
		}

		// IRL just use c.IndentedJSON(200, metric)
		// Use encoder manually
		enc := json.NewEncoder(c.Writer)
		enc.SetIndent("", "  ")
		if err := enc.Encode(metric); err != nil {
			slog.Error(err.Error())
		}
		c.Header("Content-Type", "application/json")

		c.Next()
		slog.Info("Processed OK")
	}
}
