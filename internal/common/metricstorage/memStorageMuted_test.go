package metricstorage

import (
	"fmt"
	"sort"
)

func ExampleMemStorageMuted_GetCounter() {
	ms := NewMemStorageMuted()
	ms.UpdateCounter("test_counter", 123)
	ms.UpdateCounter("test_counter_2", 456)
	keys := ms.ListCounters()
	// There is a map so output is not ordered
	sort.Strings(keys)
	fmt.Println(keys)
	// Optimistic usage
	val := ms.GetCounter("test_counter")
	fmt.Println(val)

	// Output:
	// [test_counter test_counter_2]
	// 123
}

func ExampleMemStorageMuted_GetGauge() {
	ms := NewMemStorageMuted()
	ms.UpdateGauge("test_gauge", 123.4)
	keys := ms.ListGauges()
	fmt.Println(keys)
	// Optimistic usage
	val := ms.GetGauge("test_gauge")
	fmt.Println(val)

	// Output:
	// [test_gauge]
	// 123.4
}
