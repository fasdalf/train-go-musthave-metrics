package server

type Storage interface {
	UpdateCounter(key string, value int)
	UpdateGauge(key string, value float64)
	GetCounter(key string) int
	GetGauge(key string) float64
}
type MemStorage struct {
	counters map[string]int
	gauges   map[string]float64
}

func NewMemStorage() Storage {
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
