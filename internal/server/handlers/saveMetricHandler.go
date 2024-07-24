package handlers

import (
	"encoding/json"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
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
