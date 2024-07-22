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
			if !s.HasGauge(metric.ID) {
				http.Error(c.Writer, fmt.Sprintf(`metric "%s" not found`, html.EscapeString(metric.ID)), http.StatusNotFound)
				return
			}
		case constants.CounterStr:
			if !s.HasCounter(metric.ID) {
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
		return
	}
}
