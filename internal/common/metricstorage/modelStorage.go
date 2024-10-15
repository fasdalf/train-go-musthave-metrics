package metricstorage

import (
	"context"
	"errors"
	"fmt"
	"html"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
)

type basicBatch interface {
	basicUpdater
	Commit() error
}

type basicStorage interface {
	UpdateCounter(key string, value int) error
	UpdateGauge(key string, value float64) error
	GetCounter(key string) (int, error)
	GetGauge(key string) (float64, error)
	HasCounter(key string) (bool, error)
	HasGauge(key string) (bool, error)
	ListGauges() ([]string, error)
	ListCounters() ([]string, error)
	StartBatch(ctx context.Context) (basicBatch, error)
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

type basicUpdater interface {
	UpdateCounter(key string, value int) error
	UpdateGauge(key string, value float64) error
}

func saveCommonModel(u basicUpdater, metric *apimodels.Metrics) (err error) {
	switch metric.MType {
	case constants.GaugeStr:
		if metric.Value == nil {
			return errors.New("empty metric value, float64 required")
		}
		err = u.UpdateGauge(metric.ID, *metric.Value)
	case constants.CounterStr:
		if metric.Delta == nil {
			return errors.New("empty metric delta, integer required")
		}
		delta := int(*metric.Delta)
		err = u.UpdateCounter(metric.ID, delta)
	default:
		err = fmt.Errorf(
			"invalid type \"%s\", only %s and %s supported",
			html.EscapeString(metric.MType),
			constants.GaugeStr,
			constants.CounterStr,
		)
	}

	return err
}

// SaveCommonModel save standard model or throw error
func (s *SavableModelStorage) SaveCommonModel(metric *apimodels.Metrics) (err error) {
	return saveCommonModel(s, metric)
}

// SaveCommonModels save slice of standard models or throw error
func (s *SavableModelStorage) SaveCommonModels(ctx context.Context, metrics []apimodels.Metrics) error {
	batch, err := s.StartBatch(ctx)
	if err != nil {
		return err
	}

	for _, m := range metrics {
		if err = saveCommonModel(batch, &m); err != nil {
			return err
		}
	}

	err = batch.Commit()

	return err
}
