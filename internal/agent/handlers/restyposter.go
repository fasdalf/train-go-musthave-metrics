package handlers

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/go-resty/resty/v2"
)

type restyPoster struct {
	netHTTPPoster
}

func (p *restyPoster) Post(ctx context.Context, idlog *slog.Logger, metrics []*apimodels.Metrics) error {
	body, err := compressMetrics(metrics, p.encryptionKey)
	if err != nil {
		idlog.Error("failed to prepare request body", "error", err)
		return err
	}

	client := resty.New()
	req := client.R()
	req.SetContext(ctx)
	req.SetHeader("Content-Encoding", "gzip")
	req.SetHeader("Accept-Encoding", "gzip")
	req.SetHeader("Content-Type", "application/json")

	if p.key != "" {
		hash := cryptofacade.Hash(body.Bytes(), []byte(p.key))
		req.SetHeader(constants.HeaderHashSHA256, hash)
	}

	req.SetBody(body)
	resp, err := req.Post(p.address)
	if err != nil {
		idlog.Error("send request error", "error", err)
		return fmt.Errorf("sending metrics: %w", err)
	}

	if resp != nil && resp.RawResponse.Body != nil {
		_ = resp.RawResponse.Body.Close()
	}

	if resp != nil && resp.IsError() {
		idlog.Error("response error", "error", resp.Error())
		return fmt.Errorf("sending metrics https status error: %d", resp.StatusCode())
	}
	return nil
}

func NewRestyPoster(address string, key string, encryptionKey *rsa.PublicKey) *restyPoster {
	return &restyPoster{*NewNetHTTPPoster(address, key, encryptionKey)}
}
