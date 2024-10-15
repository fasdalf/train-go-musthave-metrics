package metricstorage

import (
	"strconv"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
)

func BenchmarkModelStorageSingleUpdates(b *testing.B) {
	count := 20

	ms := NewSavableModelStorage(NewDirtyStorage(NewMemStorage()))
	m := &apimodels.Metrics{}
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			v := int64(i)
			m.MType = constants.CounterStr
			m.ID = "counter" + strconv.Itoa(i)
			m.Delta = &v
		} else {
			v := float64(i)
			m.MType = constants.GaugeStr
			m.ID = "gauge" + strconv.Itoa(i)
			m.Value = &v
		}
		_ = ms.SaveCommonModel(m)
	}
	b.ResetTimer()
	b.Run("single updates", func(b *testing.B) {
		var n int
		var f float64
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				v := float64(i)
				m.MType = constants.CounterStr
				m.ID = "counter" + strconv.Itoa(i%count)
				m.Value = &v
				_ = ms.SaveCommonModel(m)
				n, _ = ms.GetCounter("counter" + strconv.Itoa(i%count))
			} else {
				v := float64(i)
				m.MType = constants.GaugeStr
				m.ID = "gauge" + strconv.Itoa(i)
				m.Value = &v
				_ = ms.SaveCommonModel(m)
				f, _ = ms.GetGauge("gauge" + strconv.Itoa(i%count))
			}
		}
		b.ReportMetric(float64(n)/f, "useless/op")
	})
}
