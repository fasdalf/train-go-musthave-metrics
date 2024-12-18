package handlers

import (
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	s := metricstorage.NewMemStorageMuted()
	CollectMetrics(s, 0)
	assert.Equal(t, len(s.ListCounters()), 1)
	assert.Equal(t, len(s.ListGauges()), 28)
}
