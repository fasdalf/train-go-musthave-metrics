package handlers

import (
	"crypto/rsa"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewDecryptBodyHandler save metric batch to a storage
func NewDecryptBodyHandler(k *rsa.PrivateKey) func(c *gin.Context) {
	decryptBodyHandler := func(c *gin.Context) {
		if k != nil {
			// оборачиваем тело запроса в io.Reader с поддержкой декриптографирования
			cr, err := newCryptoReader(c.Request.Body, k)
			if err != nil {
				slog.Error("can't decrypt body", "error", err)
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("can't decrypt body"))
				return
			}
			defer cr.Close()
			// меняем тело запроса на новое
			c.Request.Body = cr
		}
		c.Next()
	}
	return decryptBodyHandler
}
