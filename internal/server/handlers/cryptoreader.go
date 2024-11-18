package handlers

import (
	"bytes"
	"crypto/rsa"
	"io"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"
)

// cryptoReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декриптировать получаемые от клиента данные
type cryptoReader struct {
	r io.ReadCloser
	b *bytes.Buffer
	k *rsa.PrivateKey
}

func newCryptoReader(r io.ReadCloser, k *rsa.PrivateKey) (*cryptoReader, error) {
	b1 := new(bytes.Buffer)
	_, err := b1.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	b2, err := rsacrypt.DecryptWithPrivateKey(b1.Bytes(), k)
	if err != nil {
		return nil, err
	}

	b3 := bytes.NewBuffer(b2)

	return &cryptoReader{
		r: r,
		b: b3,
		k: k,
	}, nil
}

func (c *cryptoReader) Read(p []byte) (n int, err error) {
	return c.b.Read(p)
}

func (c *cryptoReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return nil
}
