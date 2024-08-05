package metricstorage

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
