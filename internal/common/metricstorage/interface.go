package metricstorage

type Storage interface {
	UpdateCounter(key string, value int)
	UpdateGauge(key string, value float64)
	GetCounter(key string) int
	GetGauge(key string) float64
	ListGauges() []string
	ListCounters() []string
}
