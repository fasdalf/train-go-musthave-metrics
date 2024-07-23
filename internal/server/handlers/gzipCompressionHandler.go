package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GzipCompressionHandler middleware
func GzipCompressionHandler(c *gin.Context) {
	// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
	// который будем передавать следующей функции
	ow := c.Writer

	// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
	acceptEncoding := c.GetHeader("Accept-Encoding")
	supportsGzip := strings.Contains(acceptEncoding, "gzip")
	if supportsGzip {
		// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
		cw := newCompressWriter(c.Writer)
		// меняем оригинальный http.ResponseWriter на новый
		ow = cw
		cw.Header().Set("Content-Encoding", "gzip")
		// не забываем отправить клиенту все сжатые данные после завершения middleware
		defer cw.Close()
	}

	c.Writer = ow

	// проверяем, что клиент отправил серверу сжатые данные в формате gzip
	contentEncoding := c.GetHeader("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")
	if sendsGzip {
		// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
		cr, err := newCompressReader(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		// меняем тело запроса на новое
		c.Request.Body = cr
		defer cr.Close()
	}

	// передаём управление хендлеру
	c.Next()
}
