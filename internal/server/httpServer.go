package server

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"
	"log/slog"
)

func NewHTTPEngine(ms handlers.Storage) *gin.Engine {
	ginCore := gin.New()
	ginCore.Use(slogGin.New(slog.Default()))
	ginCore.Use(gin.Recovery())

	// check with and w/o trailing slash
	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET(fmt.Sprintf("/value/:%s/:%s", handlers.URLParamType, handlers.URLParamName), handlers.NewViewStatsHandler(ms))
	ginCore.POST(fmt.Sprintf("/update/:%s/:%s/:%s", handlers.URLParamType, handlers.URLParamName, handlers.URLParamValue), handlers.NewUpdateMetricHandler(ms))
	return ginCore
}
