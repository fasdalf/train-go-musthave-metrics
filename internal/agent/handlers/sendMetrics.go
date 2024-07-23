package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"

	resty "github.com/go-resty/resty/v2"
)

const URLTemplate = "http://%s/update/"

// SendMetrics sends pre collected metrics to server
func SendMetrics(s Storage, address string) {
	slog.Info("Sending metricUpdates")
	address = fmt.Sprintf(URLTemplate, address)
	client := resty.New()
	metricUpdates := make([]apimodels.Metrics, 0)
	for _, key := range s.ListCounters() {
		counter := int64(s.GetCounter(key))
		metricUpdates = append(metricUpdates, apimodels.Metrics{
			ID:    key,
			MType: constants.CounterStr,
			Delta: &counter,
			Value: nil,
		})
	}
	for _, key := range s.ListGauges() {
		gauge := s.GetGauge(key)
		metricUpdates = append(metricUpdates, apimodels.Metrics{
			ID:    key,
			MType: constants.GaugeStr,
			Delta: nil,
			Value: &gauge,
		})
	}
	for _, update := range metricUpdates {
		content, err := json.Marshal(update)
		if err != nil {
			slog.Error("error encoding request", "error", err)
			return
		}
		body := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(body)
		_, err = zb.Write(content)
		if err != nil {
			slog.Error("error compressing request", "error", err)
			return
		}

		req := client.R()
		req.SetHeader("Content-Encoding", "gzip")
		req.SetHeader("Accept-Encoding", "gzip")
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(body)
		//resp, err := req.Post(address)
		_, err = req.Post(address)
		if err != nil {
			slog.Error("Sending metrics failed", "error", err)
			continue
		}

		// TODO: Check if RAM leaks without it
		//if resp != nil && resp.RawResponse.Body != nil {
		//	_ = resp.RawResponse.Body.Close()
		//}
	}
}
