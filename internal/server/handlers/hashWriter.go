package handlers

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/gin-gonic/gin"
)

// hashWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// посчитать хеш записанного в него тела.
type hashWriter struct {
	gin.ResponseWriter
	body []byte
}

func newHashWriter(w gin.ResponseWriter) *hashWriter {
	return &hashWriter{
		ResponseWriter: w,
		body:           make([]byte, 0),
	}
}

func (c *hashWriter) Write(p []byte) (int, error) {
	c.body = append(c.body, p...)
	return len(p), nil
}

func (c *hashWriter) GetHash(key string) string {
	return cryptofacade.Hash(c.body, []byte(key))
}

func (c *hashWriter) WriteReally() (int, error) {
	return c.ResponseWriter.Write(c.body)
}
