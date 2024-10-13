package handlers

import (
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"html"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	URLParamType  = "type"
	URLParamName  = "name"
	URLParamValue = "value"
)

// NewUpdateMetricHandler update single metric value
//
// #Deprecated: for old lesson
func NewUpdateMetricHandler(s Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		mType := c.Param(URLParamType)
		mName := c.Param(URLParamName)
		mValue := c.Param(URLParamValue)
		slog.Info("Processing update", "type", mType, "name", mName, "value", mValue)

		if mName == "" {
			slog.Error("Metric not found")
			_ = c.AbortWithError(http.StatusNotFound, errors.New("metric not found"))
			return
		}

		metric := apimodels.Metrics{ID: mName, MType: mType}

		switch mType {
		case constants.GaugeStr:
			floatValue, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				slog.Error("Invalid metric value, float64 required", `value`, html.EscapeString(mType))
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid metric value"))
				return
			}
			metric.Value = &floatValue
		case constants.CounterStr:
			intValue, err := strconv.Atoi(mValue)
			if err != nil {
				slog.Error("Invalid metric value, integer required", `value`, html.EscapeString(mType))
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid metric value"))
				return
			}
			int64Value := int64(intValue)
			metric.Delta = &int64Value
		default:
			slog.Error("Type is invalid", "type", html.EscapeString(mType))
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf(
				`invalid type "%s", only "%s" and "%s" are supported`,
				html.EscapeString(mType),
				constants.GaugeStr,
				constants.CounterStr,
			))
			return
		}

		if err := s.SaveCommonModel(&metric); err != nil {
			slog.Error("can't update metric", "error", err)
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("unexpected error"))
			return
		}

		c.Next()
		slog.Info("Processed OK")
	}
}
