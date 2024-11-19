package handlers

import (
	"io"
	"strings"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/mockreaders"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"
)

func TestCryptoReaderClosedStreams(t *testing.T) {
	privK, pubK := rsacrypt.GenerateKeyPair(2048)
	_, err := newCryptoReader(&mockreaders.ErrReader{}, privK)
	_ = (&mockreaders.ErrReader{}).Close()
	if err == nil {
		t.Errorf("Expected an error on empty reader, got nil")
	}
	body1 := "not encrypted"
	rc1 := io.NopCloser(strings.NewReader(body1))
	_, err = newCryptoReader(rc1, privK)
	if err == nil {
		t.Errorf("Expected an error on unencrypted data, got nil")
	}
	crypted, err := rsacrypt.EncryptWithPublicKey([]byte(body1), pubK)
	if err != nil {
		t.Errorf("Failed to encrypt data: %v", err)
	}
	rc2 := mockreaders.NewErrCloser(crypted)
	r2, err := newCryptoReader(rc2, privK)
	if err != nil {
		t.Errorf("Failed to create reader: %v", err)
	}
	err = r2.Close()
	if err == nil {
		t.Errorf("Expected an error on close, got nil")
	}
}
