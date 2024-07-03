package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	resty "github.com/go-resty/resty/v2"
)

func SendMetrics(s metricstorage.Storage) {
	fmt.Println("Sending metrics")
	const baseURL = "http://localhost:8080/update/"
	client := resty.New()
	urls := make([]string, 0)
	for _, key := range s.ListCounters() {
		urls = append(urls, fmt.Sprintf(`%s%s/%s/%d`, baseURL, constants.CounterStr, key, s.GetCounter(key)))
	}
	for _, key := range s.ListGauges() {
		urls = append(urls, fmt.Sprintf(`%s%s/%s/%f`, baseURL, constants.GaugeStr, key, s.GetGauge(key)))
	}
	for _, url := range urls {
		resp, err := client.R().Post(url)
		if err != nil {
			fmt.Println("Error sending metrics: ", err)
			continue
		}
		if resp != nil && resp.RawResponse.Body != nil {
			_ = resp.RawResponse.Body.Close()
		}
	}
}
