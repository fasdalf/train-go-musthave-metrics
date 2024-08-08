package handlers

import (
	"encoding/json"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// SaveMetricHandler save metric to storage
func SaveMetricHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		metric := &apimodels.Metrics{}

		// IRL just use err := c.BindJSON(&metric); and c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		// use json and http on our own
		dec := json.NewDecoder(c.Request.Body)
		if err := dec.Decode(metric); err != nil {
			slog.Error("can't parse JSON", "error", err)
			http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := s.SaveCommonModel(metric)
		if err != nil {
			slog.Error("can't save metric", "error", err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Set(contextMetricResponseKey, metric)
		c.Next()
	}
}
