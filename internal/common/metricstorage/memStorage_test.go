package metricstorage

import (
	"strconv"
	"testing"
)

func BenchmarkMemStorageSingleUpdates(b *testing.B) {
	count := 20

	ms := NewMemStorage()
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			_ = ms.UpdateCounter("counter"+strconv.Itoa(i), i)
		} else {
			_ = ms.UpdateGauge("gauge"+strconv.Itoa(i), float64(i)*1.001)
		}
	}
	b.ResetTimer()
	b.Run("single updates", func(b *testing.B) {
		var n int
		var f float64
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				_ = ms.UpdateCounter("counter"+strconv.Itoa(i%count), i)
				n, _ = ms.GetCounter("counter" + strconv.Itoa(i%count))
			} else {
				_ = ms.UpdateGauge("gauge"+strconv.Itoa(i%count), float64(i)*1.001)
				f, _ = ms.GetGauge("gauge" + strconv.Itoa(i%count))
			}
		}
		b.ReportMetric(float64(n)/f, "useless/op")
	})
}
