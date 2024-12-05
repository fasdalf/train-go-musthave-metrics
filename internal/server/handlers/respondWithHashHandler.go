package handlers

import (
	"log/slog"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
)

// NewRespondWithHashHandler validate hash header and
func NewRespondWithHashHandler(key string) func(c *gin.Context) {
	validateHashHandler := func(c *gin.Context) {
		if isHashPresent, ok := c.Get(IsHashPresent); !ok || !isHashPresent.(bool) {
			slog.Info("Hash is not present in context")
			c.Next()
			return
		}

		newWriter := newHashWriter(c.Writer)
		c.Writer = newWriter
		c.Next()

		hash := newWriter.GetHash(key)
		slog.Info("setting header", "header", constants.HeaderHashSHA256, "value", hash)
		c.Header(constants.HeaderHashSHA256, hash)
		if _, err := newWriter.WriteReally(); err != nil {
			slog.Error("error writing response", "error", err)
		}
	}
	return validateHashHandler
}
