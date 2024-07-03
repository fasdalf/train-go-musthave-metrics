package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"math/rand"
	"runtime"
)

func CollectMetrics(s metricstorage.Storage) {
	fmt.Println("Collecting metrics")
	s.UpdateCounter("PollCount", s.GetCounter("PollCount")+1)
	s.UpdateGauge("RandomValue", rand.Float64())

	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	s.UpdateGauge("Alloc", float64(ms.Alloc))
	s.UpdateGauge("BuckHashSys", float64(ms.BuckHashSys))
	s.UpdateGauge("Frees", float64(ms.Frees))
	s.UpdateGauge("GCCPUFraction", ms.GCCPUFraction)
	s.UpdateGauge("GCSys", float64(ms.GCSys))
	s.UpdateGauge("HeapAlloc", float64(ms.HeapAlloc))
	s.UpdateGauge("HeapIdle", float64(ms.HeapIdle))
	s.UpdateGauge("HeapInuse", float64(ms.HeapInuse))
	s.UpdateGauge("HeapObjects", float64(ms.HeapObjects))
	s.UpdateGauge("HeapReleased", float64(ms.HeapReleased))
	s.UpdateGauge("HeapSys", float64(ms.HeapSys))
	s.UpdateGauge("LastGC", float64(ms.LastGC))
	s.UpdateGauge("Lookups", float64(ms.Lookups))
	s.UpdateGauge("MCacheInuse", float64(ms.MCacheInuse))
	s.UpdateGauge("MCacheSys", float64(ms.MCacheSys))
	s.UpdateGauge("MSpanInuse", float64(ms.MSpanInuse))
	s.UpdateGauge("MSpanSys", float64(ms.MSpanSys))
	s.UpdateGauge("Mallocs", float64(ms.Mallocs))
	s.UpdateGauge("NextGC", float64(ms.NextGC))
	s.UpdateGauge("NumForcedGC", float64(ms.NumForcedGC))
	s.UpdateGauge("NumGC", float64(ms.NumGC))
	s.UpdateGauge("OtherSys", float64(ms.OtherSys))
	s.UpdateGauge("PauseTotalNs", float64(ms.PauseTotalNs))
	s.UpdateGauge("StackInuse", float64(ms.StackInuse))
	s.UpdateGauge("StackSys", float64(ms.StackSys))
	s.UpdateGauge("Sys", float64(ms.Sys))
	s.UpdateGauge("TotalAlloc", float64(ms.TotalAlloc))
}

/*

// Can use smth. like this when all counters are in same object
//    p := Point{3, 5, "Z"}
//    pX := getAttr(&p, "X")
//
//    // Get test (int)
//    fmt.Println(pX.Int()) // 3

func getAttr(obj interface{}, fieldName string) reflect.Value {
    pointToStruct := reflect.ValueOf(obj) // addressable
    curStruct := pointToStruct.Elem()
    if curStruct.Kind() != reflect.Struct {
        panic("not struct")
    }
    curField := curStruct.FieldByName(fieldName) // type: reflect.Value
    if !curField.IsValid() {
        panic("not found:" + fieldName)
    }
    return curField
}
*/
