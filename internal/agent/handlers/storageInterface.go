package handlers

type Storage interface {
	UpdateCounter(key string, value int)
	UpdateGauge(key string, value float64)
	GetCounter(key string) int
	GetGauge(key string) float64
	// Agent don't use these methods.
	// HasCounter(key string) bool
	// HasGauge(key string) bool
	ListGauges() []string
	ListCounters() []string
}
