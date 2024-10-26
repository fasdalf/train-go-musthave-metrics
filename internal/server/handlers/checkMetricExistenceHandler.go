package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
)

// CheckMetricExistenceHandler checks metrics has value set
// Previous handler should add valid metric to context
func CheckMetricExistenceHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		metric := &apimodels.Metrics{}

		// IRL just use err := c.BindJSON(&metric); and c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		// use json and http on our own
		dec := json.NewDecoder(c.Request.Body)
		if err := dec.Decode(metric); err != nil {
			slog.Error("can't parse JSON", "error", err)
			// TODO: replace with c.AbortWithError to skip later middlewares
			http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case constants.GaugeStr:
			if h, err := s.HasGauge(metric.ID); err != nil || !h {
				if err != nil {
					slog.Error("can't check metric existence", "id", metric.ID, "error", err)
					http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
					return
				}
				slog.Error("metric not found", "id", metric.ID)
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(metric.ID)), http.StatusNotFound)
				return
			}
		case constants.CounterStr:
			if h, err := s.HasCounter(metric.ID); err != nil || !h {
				if err != nil {
					slog.Error("can't check metric existence", "id", metric.ID, "error", err)
					http.Error(c.Writer, `unexpected error`, http.StatusInternalServerError)
					return
				}
				slog.Error("metric not found", "id", metric.ID)
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(metric.ID)), http.StatusNotFound)
				return
			}
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
