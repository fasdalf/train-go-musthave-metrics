package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/go-resty/resty/v2"
	"log/slog"
)

type restyPoster struct{}

func (p *restyPoster) Post(ctx context.Context, idlog *slog.Logger, body *bytes.Buffer, key string, address string) error {
	client := resty.New()
	req := client.R()
	req.SetContext(ctx)
	req.SetHeader("Content-Encoding", "gzip")
	req.SetHeader("Accept-Encoding", "gzip")
	req.SetHeader("Content-Type", "application/json")

	if key != "" {
		hash := cryptofacade.Hash(body.Bytes(), []byte(key))
		req.SetHeader(constants.HashSHA256, hash)
	}

	req.SetBody(body)
	resp, err := req.Post(address)
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

func NewRestyPoster() *restyPoster {
	return &restyPoster{}
}
