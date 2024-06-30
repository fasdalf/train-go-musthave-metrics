package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/constants"
	"net/http"
)

func SendMetrics(s metricstorage.Storage) {
	fmt.Println("Sending metrics")
	const baseURL = "http://localhost:8080/update/"
	for _, key := range s.ListCounters() {
		url := fmt.Sprintf(`%s%s/%s/%d`, baseURL, constants.CounterStr, key, s.GetCounter(key))
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println("Error sending metrics: ", err)
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
	for _, key := range s.ListGauges() {
		url := fmt.Sprintf(`%s%s/%s/%f`, baseURL, constants.GaugeStr, key, s.GetGauge(key))
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println("Error sending metrics: ", err)
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}
