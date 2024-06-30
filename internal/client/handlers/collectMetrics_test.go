package handlers

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	s := metricstorage.NewMemStorage()
	CollectMetrics(s)
	assert.Equal(t, len(s.ListCounters()), 1)
	assert.Equal(t, len(s.ListGauges()), 28)
}