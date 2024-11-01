package metricstorage

import (
	"fmt"
)

func ExampleMemStorageMuted_GetCounter() {
	ms := NewMemStorageMuted()
	ms.UpdateCounter("test_counter", 123)
	ms.UpdateCounter("test_counter_2", 456)
	// Optimistic usage
	val := ms.GetCounter("test_counter")
	fmt.Println(val)
	// Output:
	// 123
}

func ExampleMemStorageMuted_GetGauge() {
	ms := NewMemStorageMuted()
	ms.UpdateGauge("test_gauge", 123.4)
	// Optimistic usage
	val := ms.GetGauge("test_gauge")
	fmt.Println(val)

	// Output:
	// 123.4
}
