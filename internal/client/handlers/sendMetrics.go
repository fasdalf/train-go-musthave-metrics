package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/constants"
	"net/http"
)

func SendMetrics(s metricstorage.Storage) {
	fmt.Println("Sending metrics")
	const baseUrl = "http://127.0.0.1:8080/update/"
	for _, key := range s.ListCounters() {
		url := fmt.Sprintf(`%s%s/%s/%d`, baseUrl, constants.CounterStr, key, s.GetCounter(key))
		_, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println("Error sending metrics: ", url, err)
		}
	}
	for _, key := range s.ListGauges() {
		url := fmt.Sprintf(`%s%s/%s/%f`, baseUrl, constants.GaugeStr, key, s.GetGauge(key))
		_, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println("Error sending metrics: ", url, err)
		}
	}
}
