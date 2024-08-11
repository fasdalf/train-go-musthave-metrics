package metricstorage

import (
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"html"
	"log/slog"
)

type MemStorage struct {
	counters map[string]int
	gauges   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counters: make(map[string]int),
		gauges:   make(map[string]float64),
	}
}

func (s *MemStorage) UpdateCounter(key string, value int) {
	s.counters[key] += value
}

func (s *MemStorage) UpdateGauge(key string, value float64) {
	s.gauges[key] = value
}

func (s *MemStorage) GetCounter(key string) int {
	return s.counters[key]
}

func (s *MemStorage) GetGauge(key string) float64 {
	return s.gauges[key]
}

func (s *MemStorage) HasCounter(key string) bool {
	_, ok := s.counters[key]
	return ok
}

func (s *MemStorage) HasGauge(key string) bool {
	_, ok := s.gauges[key]
	return ok
}

func (s *MemStorage) ListGauges() []string {
	keys := make([]string, 0, len(s.gauges))
	for k := range s.gauges {
		keys = append(keys, k)
	}
	return keys
}

func (s *MemStorage) ListCounters() []string {
	keys := make([]string, 0, len(s.counters))
	for k := range s.counters {
		keys = append(keys, k)
	}
	return keys
}

// MemStorageWithSave Just to try go's inheritance
type MemStorageWithSave struct {
	MemStorage
}

func NewMemStorageWithSave() *MemStorageWithSave {
	return &MemStorageWithSave{
		*NewMemStorage(),
	}
}

// SaveCommonModel save standard model or throw error
func (s *MemStorageWithSave) SaveCommonModel(metric *apimodels.Metrics) error {
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
