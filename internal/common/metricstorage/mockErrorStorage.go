package metricstorage

import (
	"io"
)

type MockErrorStorage struct {
	MemStorage
	WithError bool
}

func (s *MockErrorStorage) HasCounter(key string) (bool, error) {
	if s.WithError {
		return false, io.EOF
	}
	return s.MemStorage.HasCounter(key)
}

func (s *MockErrorStorage) HasGauge(key string) (bool, error) {
	if s.WithError {
		return false, io.EOF
	}
	return s.MemStorage.HasGauge(key)
}

func (s *MockErrorStorage) UpdateCounter(key string, value int) error {
	if s.WithError {
		return io.EOF
	}
	return s.MemStorage.UpdateCounter(key, value)
}

func (s *MockErrorStorage) UpdateGauge(key string, value float64) error {
	if s.WithError {
		return io.EOF
	}
	return s.MemStorage.UpdateGauge(key, value)
}

func (s *MockErrorStorage) GetCounter(key string) (int, error) {
	if s.WithError {
		return 0, io.EOF
	}
	return s.MemStorage.GetCounter(key)
}

func (s *MockErrorStorage) GetGauge(key string) (float64, error) {
	if s.WithError {
		return 0, io.EOF
	}
	return s.MemStorage.GetGauge(key)
}
