package server

import (
	"fmt"
	"log/slog"

	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"

	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"
)

func NewHTTPEngine(ms handlers.Storage) *gin.Engine {
	ginCore := gin.New()
	ginCore.Use(slogGin.New(slog.Default()))
	ginCore.Use(gin.Recovery())
	// import "github.com/gin-contrib/gzip"
	//ginCore.Use(gzip.Gzip(gzip.DefaultCompression))
	// uses internal/server/handlers/gzip.go + internal/server/handlers/gzipCompressionHandler.go
	ginCore.Use(handlers.GzipCompressionHandler)

	// check with and w/o trailing slash
	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET(fmt.Sprintf("/value/:%s/:%s", handlers.URLParamType, handlers.URLParamName), handlers.NewViewStatsHandler(ms))
	ginCore.POST("/value/", handlers.CheckMetricExistenceHandler(ms), handlers.MetricValueResponseHandler(ms))
	ginCore.POST(fmt.Sprintf("/update/:%s/:%s/:%s", handlers.URLParamType, handlers.URLParamName, handlers.URLParamValue), handlers.NewUpdateMetricHandler(ms))
	ginCore.POST("/update/", handlers.SaveMetricHandler(ms), handlers.MetricValueResponseHandler(ms))
	return ginCore
}
