package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	resty "github.com/go-resty/resty/v2"
	"log/slog"
)

const URLTemplate = "http://%s/update/"

var address = ``

func SendMetrics(s metricstorage.Storage, address string) {
	slog.Info("Sending metrics")
	address = fmt.Sprintf(URLTemplate, address)
	client := resty.New()
	urls := make([]string, 0)
	for _, key := range s.ListCounters() {
		urls = append(urls, fmt.Sprintf(`%s%s/%s/%d`, address, constants.CounterStr, key, s.GetCounter(key)))
	}
	for _, key := range s.ListGauges() {
		urls = append(urls, fmt.Sprintf(`%s%s/%s/%f`, address, constants.GaugeStr, key, s.GetGauge(key)))
	}
	for _, url := range urls {
		resp, err := client.R().Post(url)
		if err != nil {
			slog.Error(`Sending metrics failed`, `error`, err)
			continue
		}
		if resp != nil && resp.RawResponse.Body != nil {
			_ = resp.RawResponse.Body.Close()
		}
	}
}
