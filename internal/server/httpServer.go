package server

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"net/http"
)

func NewHTTPEngine(ms metricstorage.Storage) *gin.Engine {
	ginCore := gin.New()
	ginCore.Use(ginlogrus.Logger(logrus.New()))
	ginCore.Use(gin.Recovery())

	// check with and w/o trailing slash
	ginCore.GET("/", gin.WrapF(handlers.NewIndexHandler(ms)))
	ginCore.GET("/value/:type/:name", handlers.NewViewStatsHandler(ms))
	ginCore.POST("/update/:type/:name/:value", gin.WrapH(http.StripPrefix("/update/", handlers.NewUpdateMetricHandler(ms))))
	return ginCore
}
