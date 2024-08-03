package server

import (
	"fmt"
	"log/slog"

	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"

	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"
)

func NewRoutingEngine(ms handlers.Storage, fs handlers.FileStorage) *gin.Engine {
	ginCore := gin.New()
	// IRL just use ginCore.Use(slogGin.New(slog.Default())) from slogGin "github.com/samber/slog-gin"
	// "We have it at home" logging. Uses ##@@
	// TODO: WIP ##@@
	ginCore.Use(slogGin.New(slog.Default()))
	ginCore.Use(gin.Recovery())
	// IRL just use ginCore.Use(gzip.Gzip(gzip.DefaultCompression)) from "github.com/gin-contrib/gzip"
	// "We have it at home" compression. Uses internal/server/handlers/gzip.go + internal/server/handlers/gzipCompressionHandler.go
	ginCore.Use(handlers.GzipCompressionHandler)

	// check with and w/o trailing slash
	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET(fmt.Sprintf("/value/:%s/:%s", handlers.URLParamType, handlers.URLParamName), handlers.NewViewStatsHandler(ms))
	ginCore.POST("/value/", handlers.CheckMetricExistenceHandler(ms), handlers.MetricValueResponseHandler(ms))
	ginCore.POST(fmt.Sprintf("/update/:%s/:%s/:%s", handlers.URLParamType, handlers.URLParamName, handlers.URLParamValue), handlers.NewUpdateMetricHandler(ms))
	updatePipeline := []gin.HandlerFunc{handlers.SaveMetricHandler(ms), handlers.MetricValueResponseHandler(ms)}
	if fs != nil {
		updatePipeline = append(updatePipeline, handlers.NewSaveToFileHandler(fs))
	}
	ginCore.POST("/update/", updatePipeline...)
	return ginCore
}
