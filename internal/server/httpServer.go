package server

import (
	"crypto/rsa"
	"fmt"

	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func NewRoutingEngine(ms handlers.Storage, db handlers.Pingable, retryer handlers.Retryer, key string, decryptionKey *rsa.PrivateKey) *gin.Engine {
	ginCore := gin.New()
	ginCore.RedirectTrailingSlash = false
	ginCore.RedirectFixedPath = false
	ginCore.Use(gin.Recovery())
	// IRL just use ginCore.Use(slogGin.New(slog.Default())) from slogGin "github.com/samber/slog-gin"
	// "We have it at home" logging. Uses .../loggingResponseWriter.go and .../slogHandler.go
	ginCore.Use(handlers.SlogHandler)
	if key != "" {
		ginCore.Use(handlers.NewValidateHashHandler(key))
		ginCore.Use(handlers.NewRespondWithHashHandler(key))
	}
	// IRL just use ginCore.Use(gzip.Gzip(gzip.DefaultCompression)) from "github.com/gin-contrib/gzip"
	// "We have it at home" compression. Uses .../compressWriter.go and .../gzipCompressionHandler.go
	ginCore.Use(handlers.GzipCompressionHandler)
	// IRL use regular TLS.
	// "We have it at home" RSA.
	ginCore.Use(handlers.NewDecryptBodyHandler(decryptionKey))

	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET("/ping", handlers.NewPingDBHandler(db))
	ginCore.GET(fmt.Sprintf("/value/:%s/:%s", handlers.URLParamType, handlers.URLParamName), handlers.NewViewStatsHandler(ms))

	valuePipeline := []gin.HandlerFunc{handlers.CheckMetricExistenceHandler(ms), handlers.MetricValueResponseHandler(ms)}
	// Path with slash is for broken test in iteration 14
	ginCore.POST("/value/", valuePipeline...)
	ginCore.POST("/value", valuePipeline...)
	ginCore.POST(fmt.Sprintf("/update/:%s/:%s/:%s", handlers.URLParamType, handlers.URLParamName, handlers.URLParamValue), handlers.NewUpdateMetricHandler(ms))
	ginCore.POST("/update/", handlers.SaveMetricHandler(ms), handlers.MetricValueResponseHandler(ms))
	ginCore.POST("/updates/", handlers.SaveMetricsHandler(ms, retryer))
	return ginCore
}
