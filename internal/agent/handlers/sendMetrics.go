package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"

	resty "github.com/go-resty/resty/v2"
)

const URLTemplate = "http://%s/updates/"

type Retryer interface {
	Try(ctx context.Context, do func() error, isRetryable func(err error) bool) (int, error)
}

var ErrTransport = errors.New("resty error")

// SendMetrics sends pre collected metrics to server
func SendMetrics(ctx context.Context, s Storage, address string, r Retryer) {
	slog.Info("Sending metricUpdates")
	address = fmt.Sprintf(URLTemplate, address)

	doJob := func() error {
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

		if len(metricUpdates) == 0 {
			return nil
		}

		content, err := json.Marshal(metricUpdates)
		if err != nil {
			return errors.Join(fmt.Errorf("encoding request: %w", err), ErrTransport)
		}
		body := new(bytes.Buffer)
		zb := gzip.NewWriter(body)
		_, err = zb.Write(content)
		if err != nil {
			return fmt.Errorf("compressing request: %w", err)
		}
		_ = zb.Close()

		req := client.R()
		req.SetHeader("Content-Encoding", "gzip")
		req.SetHeader("Accept-Encoding", "gzip")
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(body)
		//resp, err := req.Post(address)
		resp, err := req.Post(address)
		if err != nil {
			return fmt.Errorf("sending metrics: %w", err)
		}

		if resp != nil && resp.RawResponse.Body != nil {
			_ = resp.RawResponse.Body.Close()
		}

		if resp != nil && resp.IsError() {
			return fmt.Errorf("sending metrics https status error: %d", resp.StatusCode())
		}

		return nil
	}

	isRecoverable := func(err error) bool {
		var netErr net.Error = nil
		if errors.As(err, &netErr) && netErr != nil {
			return true
		}
		return false
	}

	if _, err := r.Try(ctx, doJob, isRecoverable); err != nil {
		slog.Info("SendMetrics error", "error", err)
	}
}
