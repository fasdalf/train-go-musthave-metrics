package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"log/slog"
	"net/http"
)

type netHTTPPoster struct{}

func (p *netHTTPPoster) Post(ctx context.Context, idlog *slog.Logger, body *bytes.Buffer, key string, address string) error {
	request, err := http.NewRequest(http.MethodPost, address, body)
	if err != nil {
		idlog.Error("init request error", "error", err)
		return fmt.Errorf("sending metrics: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")

	if key != "" {
		hash := cryptofacade.Hash(body.Bytes(), []byte(key))
		request.Header.Set(constants.HashSHA256, hash)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		idlog.Error("send request error", "error", err)
		return fmt.Errorf("sending metrics: %w", err)
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp != nil && resp.StatusCode != http.StatusOK {
		idlog.Error("response error", "error", resp.StatusCode)
		return fmt.Errorf("sending metrics https status error: %d", resp.StatusCode)
	}
	return nil
}

func NewNetHTTPPoster() *netHTTPPoster {
	return &netHTTPPoster{}
}
