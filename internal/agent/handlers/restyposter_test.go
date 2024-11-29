package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"
)

func TestRestyPoster_Post_Success(t *testing.T) {
	// TODO: ##@@ add more test cases
	wantStatus := http.StatusOK
	_, pub := rsacrypt.GenerateKeyPair(2048)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(wantStatus)
	}))
	// останавливаем сервер после завершения теста
	defer srv.Close()

	url := strings.Replace(srv.URL, "http://", "", 1)
	fmt.Println(url)
	poster := NewRestyPoster(url, "key", pub)

	i64 := int64(10)
	update := []*apimodels.Metrics{{
		ID:    "wow",
		MType: constants.CounterStr,
		Delta: &i64,
		Value: nil,
	}}
	err := poster.Post(context.Background(), slog.Default(), update)
	if err != nil {
		t.Fatalf("error posting metrics: %v", err)
	}
}
