package metricstorage

import "testing"

func TestMockErrorStorage(t *testing.T) {
	mms := &MockErrorStorage{
		MemStorage: *NewMemStorage(),
		WithError:  true,
	}
	ms := NewSavableModelStorage(mms)

	ms.UpdateCounter("counter1", 10)
	ms.HasCounter("counter1")
	ms.GetCounter("counter1")
	ms.UpdateGauge("Gauge", 10)
	ms.HasGauge("Gauge")
	ms.GetGauge("Gauge")

	mms.WithError = false

	ms.UpdateCounter("counter1", 10)
	ms.HasCounter("counter1")
	ms.GetCounter("counter1")
	ms.UpdateGauge("Gauge", 10)
	ms.HasGauge("Gauge")
	ms.GetGauge("Gauge")
}
