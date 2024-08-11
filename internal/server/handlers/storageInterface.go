package handlers

import "github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"

type Storage interface {
	UpdateCounter(key string, value int)
	UpdateGauge(key string, value float64)
	GetCounter(key string) int
	GetGauge(key string) float64
	HasCounter(key string) bool
	HasGauge(key string) bool
	ListGauges() []string
	ListCounters() []string
	SaveCommonModel(metric *apimodels.Metrics) error
}

type FileStorage interface {
	SaveWithInterval() error
}
