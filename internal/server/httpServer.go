package server

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"
	"log/slog"
)

func NewHTTPEngine(ms metricstorage.Storage) *gin.Engine {
	ginCore := gin.New()
	ginCore.Use(slogGin.New(slog.Default()))
	ginCore.Use(gin.Recovery())

	// check with and w/o trailing slash
	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET("/value/:type/:name", handlers.NewViewStatsHandler(ms))
	ginCore.POST("/update/:type/:name/:value", handlers.NewUpdateMetricHandler(ms))
	return ginCore
}
