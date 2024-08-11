package server

import (
	"database/sql"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/jsonofflinestorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func NewRoutingEngine(ms handlers.Storage, saved jsonofflinestorage.SavedChan, db *sql.DB, retryer handlers.Retryer) *gin.Engine {
	ginCore := gin.New()
	ginCore.Use(gin.Recovery())
	// IRL just use ginCore.Use(slogGin.New(slog.Default())) from slogGin "github.com/samber/slog-gin"
	// "We have it at home" logging. Uses .../loggingResponseWriter.go and .../slogHandler.go
	ginCore.Use(handlers.SlogHandler)
	// IRL just use ginCore.Use(gzip.Gzip(gzip.DefaultCompression)) from "github.com/gin-contrib/gzip"
	// "We have it at home" compression. Uses .../compressWriter.go and .../gzipCompressionHandler.go
	ginCore.Use(handlers.GzipCompressionHandler)

	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET("/ping", handlers.NewPingDBHandler(db))
	ginCore.GET(fmt.Sprintf("/value/:%s/:%s", handlers.URLParamType, handlers.URLParamName), handlers.NewViewStatsHandler(ms))
	ginCore.POST("/value", handlers.CheckMetricExistenceHandler(ms), handlers.MetricValueResponseHandler(ms))
	ginCore.POST(fmt.Sprintf("/update/:%s/:%s/:%s", handlers.URLParamType, handlers.URLParamName, handlers.URLParamValue), handlers.NewUpdateMetricHandler(ms))
	updatePipeline := []gin.HandlerFunc{handlers.SaveMetricHandler(ms), handlers.MetricValueResponseHandler(ms)}
	if saved != nil {
		updatePipeline = append(updatePipeline, handlers.NewSaveToFileHandler(saved))
	}
	ginCore.POST("/update/", updatePipeline...)
	updatesPipeline := []gin.HandlerFunc{handlers.SaveMetricsHandler(ms, retryer)}
	if saved != nil {
		updatesPipeline = append(updatesPipeline, handlers.NewSaveToFileHandler(saved))
	}
	ginCore.POST("/updates/", updatesPipeline...)
	return ginCore
}
