package handlers

import (
	"context"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
)

type Storage interface {
	UpdateCounter(key string, value int) error
	UpdateGauge(key string, value float64) error
	GetCounter(key string) (int, error)
	GetGauge(key string) (float64, error)
	HasCounter(key string) (bool, error)
	HasGauge(key string) (bool, error)
	ListGauges() ([]string, error)
	ListCounters() ([]string, error)
	SaveCommonModel(metric *apimodels.Metrics) error
	SaveCommonModels(ctx context.Context, metrics []apimodels.Metrics) error
}

type FileStorage interface {
	SaveWithInterval() error
}
