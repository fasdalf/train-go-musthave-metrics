package handlers

import (
	"bytes"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
)

const (
	IsHashPresent = "IsHashPresent"
)

// NewValidateHashHandler save metric batch to a storage
func NewValidateHashHandler(key string) func(c *gin.Context) {
	validateHashHandler := func(c *gin.Context) {
		requestHash := c.GetHeader(constants.HashSHA256)

		if requestHash == "" {
			c.Set(IsHashPresent, false)
			c.Next()
			return
		}

		ByteBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(ByteBody))
		realHash := cryptofacade.Hash(ByteBody, []byte(key))

		if realHash != requestHash {
			slog.Error("header value is invalid", "realHash", realHash, "requestHash", requestHash)
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s header value is invalid: %s", constants.HashSHA256, requestHash))
			return
		}

		c.Set(IsHashPresent, true)
		c.Next()
		return
	}
	return validateHashHandler
}
