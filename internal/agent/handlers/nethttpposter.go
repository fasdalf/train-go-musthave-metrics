package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/localip"
)

const URLTemplate = "http://%s/updates/"

type netHTTPPoster struct {
	address       string
	key           string
	encryptionKey *rsa.PublicKey
}

func (p *netHTTPPoster) Post(ctx context.Context, idlog *slog.Logger, metrics []*apimodels.Metrics) error {
	body, err := compressMetrics(metrics, p.encryptionKey)
	if err != nil {
		idlog.Error("failed to prepare request body", "error", err)
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.address, body)
	if err != nil {
		idlog.Error("init request error", "error", err)
		return fmt.Errorf("sending metrics: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set(constants.HeaderRealIP, localip.GetLocalIP().String())

	if p.key != "" {
		hash := cryptofacade.Hash(body.Bytes(), []byte(p.key))
		request.Header.Set(constants.HeaderHashSHA256, hash)
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

// compressMetrics compresses the metrics using gzip.
func compressMetrics(metricUpdates []*apimodels.Metrics, encryptionKey *rsa.PublicKey) (*bytes.Buffer, error) {
	content, err := json.Marshal(metricUpdates)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("encoding request: %w", err), ErrTransport)
	}
	body := new(bytes.Buffer)
	zb := gzip.NewWriter(body)
	defer zb.Close()
	if encryptionKey != nil {
		content, err = rsacrypt.EncryptWithPublicKey(content, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("encrypting request: %w", err)
		}
	}
	_, err = zb.Write(content)
	if err != nil {
		return nil, fmt.Errorf("compressing request: %w", err)
	}

	return body, nil
}

func NewNetHTTPPoster(address string, key string, encryptionKey *rsa.PublicKey) *netHTTPPoster {
	address = fmt.Sprintf(URLTemplate, address)
	return &netHTTPPoster{
		address:       address,
		key:           key,
		encryptionKey: encryptionKey,
	}
}
