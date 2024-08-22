package metricstorage

type MemStorageMuted struct {
	MemStorage
}

func NewMemStorageMuted() *MemStorageMuted {
	return &MemStorageMuted{
		MemStorage: *NewMemStorage(),
	}
}

func (s *MemStorageMuted) UpdateCounter(key string, value int) {
	_ = s.MemStorage.UpdateCounter(key, value)
}

func (s *MemStorageMuted) UpdateGauge(key string, value float64) {
	_ = s.MemStorage.UpdateGauge(key, value)
}

func (s *MemStorageMuted) GetCounter(key string) int {
	r, _ := s.MemStorage.GetCounter(key)
	return r
}

func (s *MemStorageMuted) GetGauge(key string) float64 {
	r, _ := s.MemStorage.GetGauge(key)
	return r
}

func (s *MemStorageMuted) ListGauges() []string {
	keys, _ := s.MemStorage.ListGauges()
	return keys
}

func (s *MemStorageMuted) ListCounters() []string {
	keys, _ := s.MemStorage.ListCounters()
	return keys
}
