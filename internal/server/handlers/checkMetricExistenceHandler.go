package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"html"
	"log/slog"
	"net/http"
)

// CheckMetricExistenceHandler checks metrics has value set
// Previous handler should add valid metric to context
func CheckMetricExistenceHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		slog.Info("#trace CheckMetricExistenceHandler")
		metric := &apimodels.Metrics{}

		// IRL just use err := c.BindJSON(&metric); and c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		// use json and http on our own
		slog.Info("#trace CheckMetricExistenceHandler")
		dec := json.NewDecoder(c.Request.Body)
		slog.Info("#trace CheckMetricExistenceHandler")
		if err := dec.Decode(metric); err != nil {
			slog.Info("#trace CheckMetricExistenceHandler")
			slog.Error("can't parse JSON", "error", err)
			http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
			slog.Info("#trace CheckMetricExistenceHandler")
			return
		}
		slog.Info("#trace CheckMetricExistenceHandler")

		slog.Info("#trace CheckMetricExistenceHandler")
		switch metric.MType {
		case constants.GaugeStr:
			slog.Info("#trace CheckMetricExistenceHandler")
			if h, err := s.HasGauge(metric.ID); err != nil || !h {
				slog.Info("#trace CheckMetricExistenceHandler")
				if err != nil {
					slog.Info("#trace CheckMetricExistenceHandler")
					slog.Error("can't check metric existence", "id", metric.ID, "error", err)
					http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
					slog.Info("#trace CheckMetricExistenceHandler")
					return
				}
				slog.Info("#trace CheckMetricExistenceHandler")
				slog.Error("metric not found", "id", metric.ID)
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(metric.ID)), http.StatusNotFound)
				slog.Info("#trace CheckMetricExistenceHandler")
				return
			}
		case constants.CounterStr:
			slog.Info("#trace CheckMetricExistenceHandler")
			if h, err := s.HasCounter(metric.ID); err != nil || !h {
				slog.Info("#trace CheckMetricExistenceHandler")
				if err != nil {
					slog.Info("#trace CheckMetricExistenceHandler")
					slog.Error("can't check metric existence", "id", metric.ID, "error", err)
					http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
					slog.Info("#trace CheckMetricExistenceHandler")
					return
				}
				slog.Info("#trace CheckMetricExistenceHandler")
				slog.Error("metric not found", "id", metric.ID)
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(metric.ID)), http.StatusNotFound)
				slog.Info("#trace CheckMetricExistenceHandler")
				return
			}
		default:
			slog.Info("#trace CheckMetricExistenceHandler")
			slog.Error("Invalid type", "type", metric.MType)
			http.Error(c.Writer, fmt.Sprintf(
				"Invalid type, only %s and %s supported",
				constants.GaugeStr,
				constants.CounterStr,
			), http.StatusBadRequest)
			slog.Info("#trace CheckMetricExistenceHandler")
			return
		}

		slog.Info("#trace CheckMetricExistenceHandler")
		slog.Info("setting to context", "key", contextMetricResponseKey, "value", metric)
		c.Set(contextMetricResponseKey, metric)
		slog.Info("#trace CheckMetricExistenceHandler")
		c.Next()
		slog.Info("#trace CheckMetricExistenceHandler")
	}
}
