package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Pingable interface {
	Ping() error
}

// NewPingDBHandler create handler to check DB connection
func NewPingDBHandler(db Pingable) func(c *gin.Context) {
	pingDBHandler := func(c *gin.Context) {
		if db == nil {
			slog.Error("DATABASE_DSN is empty")
			http.Error(c.Writer, "DATABASE_DSN is empty", http.StatusBadRequest)
			return
		}
		if err := db.Ping(); err != nil {
			slog.Error("DB connection lost", "error", err)
			http.Error(c.Writer, "DB connection lost", http.StatusBadRequest)
			return
		}

		c.Header("Content-Type", "text/plain")
		if _, err := c.Writer.Write([]byte("OK")); err != nil {
			slog.Error("response write error", "error", err)
		}

		slog.Info("Processed OK")
	}

	return pingDBHandler
}
