package metricstorage

import (
	"context"
	"sync"
)

type MemStorage struct {
	counters    map[string]int
	countersRWM *sync.RWMutex
	gauges      map[string]float64
	gaugesRWM   *sync.RWMutex
	inBatch     bool
}

type MemBatch struct {
	ms *MemStorage
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counters:    make(map[string]int),
		countersRWM: new(sync.RWMutex),
		gauges:      make(map[string]float64),
		gaugesRWM:   new(sync.RWMutex),
	}
}

func (s *MemStorage) UpdateCounter(key string, value int) error {
	if !s.inBatch {
		s.countersRWM.Lock()
		defer s.countersRWM.Unlock()
	}
	s.counters[key] += value
	return nil
}

func (s *MemStorage) UpdateGauge(key string, value float64) error {
	if !s.inBatch {
		s.gaugesRWM.Lock()
		defer s.gaugesRWM.Unlock()
	}
	s.gauges[key] = value
	return nil
}

func (s *MemStorage) GetCounter(key string) (int, error) {
	s.countersRWM.RLock()
	defer s.countersRWM.RUnlock()
	return s.counters[key], nil
}

func (s *MemStorage) GetGauge(key string) (float64, error) {
	s.gaugesRWM.RLock()
	defer s.gaugesRWM.RUnlock()
	return s.gauges[key], nil
}

func (s *MemStorage) HasCounter(key string) (bool, error) {
	s.countersRWM.RLock()
	defer s.countersRWM.RUnlock()
	_, ok := s.counters[key]
	return ok, nil
}

func (s *MemStorage) HasGauge(key string) (bool, error) {
	s.gaugesRWM.RLock()
	defer s.gaugesRWM.RUnlock()
	_, ok := s.gauges[key]
	return ok, nil
}

func (s *MemStorage) ListGauges() ([]string, error) {
	s.gaugesRWM.RLock()
	defer s.gaugesRWM.RUnlock()
	keys := make([]string, 0, len(s.gauges))
	for k := range s.gauges {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *MemStorage) ListCounters() ([]string, error) {
	s.countersRWM.RLock()
	defer s.countersRWM.RUnlock()
	keys := make([]string, 0, len(s.counters))
	for k := range s.counters {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *MemStorage) StartBatch(ctx context.Context) (basicBatch, error) {
	s.countersRWM.Lock()
	s.gaugesRWM.Lock()
	return &MemBatch{ms: s}, nil
}

func (b *MemBatch) UpdateCounter(key string, value int) error {
	return b.ms.UpdateCounter(key, value)
}
func (b *MemBatch) UpdateGauge(key string, value float64) error {
	return b.ms.UpdateGauge(key, value)

}
func (b *MemBatch) Commit() error {
	b.ms.inBatch = false
	b.ms.countersRWM.Unlock()
	b.ms.gaugesRWM.Unlock()
	return nil
}
