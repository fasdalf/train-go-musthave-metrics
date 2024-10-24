package metricstorage

import (
	"fmt"
	"strconv"
	"testing"
)

func ExampleMemStorage_GetCounter() {
	ms := NewMemStorage()
	ms.UpdateCounter("test_counter", 123)
	ms.UpdateCounter("test_counter_2", 456)
	// Optimistic usage
	val, _ := ms.GetCounter("test_counter")
	fmt.Println(val)
	// Expected usage
	val, err := ms.GetCounter("test_counter_2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}
	// Double-check before use or when value not needed
	has, err := ms.HasCounter("test_counter_3")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(has)
	}
	// Output:
	// 123
	// 456
	// false
}

func ExampleMemStorage_GetGauge() {
	ms := NewMemStorage()
	ms.UpdateGauge("test_gauge", 123.4)
	ms.UpdateGauge("test_gauge_2", 456.7)
	// Optimistic usage
	val, _ := ms.GetGauge("test_gauge")
	fmt.Println(val)
	// Expected usage
	val, err := ms.GetGauge("test_gauge_2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}
	// Double-check before use or when value not needed
	has, err := ms.HasGauge("test_gauge_3")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(has)
	}
	// Output:
	// 123.4
	// 456.7
	// false
}

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
