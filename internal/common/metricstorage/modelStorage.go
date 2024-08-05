package metricstorage

import (
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"html"
	"log/slog"
)

type basicStorage interface {
	UpdateCounter(key string, value int)
	UpdateGauge(key string, value float64)
	GetCounter(key string) int
	GetGauge(key string) float64
	HasCounter(key string) bool
	HasGauge(key string) bool
	ListGauges() []string
	ListCounters() []string
}

// SavableModelStorage wrapper around storage interface
type SavableModelStorage struct {
	basicStorage
}

func NewSavableModelStorage(bs basicStorage) *SavableModelStorage {
	return &SavableModelStorage{
		basicStorage: bs,
	}
}

// SaveCommonModel save standard model or throw error
func (s *SavableModelStorage) SaveCommonModel(metric *apimodels.Metrics) error {
	switch metric.MType {
	case constants.GaugeStr:
		if metric.Value == nil {
			return errors.New("empty metric value, float64 required")
		}
		s.UpdateGauge(metric.ID, *metric.Value)
		slog.Info("gauge value set", "key", metric.ID, "new", s.GetGauge(metric.ID))
	case constants.CounterStr:
		if metric.Delta == nil {
			return errors.New("empty metric delta, integer required")
		}
		delta := int(*metric.Delta)
		s.UpdateCounter(metric.ID, delta)
		slog.Info("counter value set", "key", metric.ID, "new", s.GetCounter(metric.ID))
	default:
		return fmt.Errorf(
			"invalid type \"%s\", only %s and %s supported",
			html.EscapeString(metric.MType),
			constants.GaugeStr,
			constants.CounterStr,
		)
	}

	return nil
}
