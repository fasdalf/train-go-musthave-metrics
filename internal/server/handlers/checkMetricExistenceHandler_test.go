package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/gin-gonic/gin"
)

// CheckMetricExistenceHandler checks metrics has value set
func TestCheckMetricExistenceHandler(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		withError   bool
		wantErr     bool
	}{
		{
			name:        "1",
			requestBody: "Hello, world!",
			wantErr:     true,
		},
		{
			name:        "counter, ok",
			requestBody: `{"id": "title1", "type": "counter"}`,
			wantErr:     false,
		},
		{
			name:        "gauge, ok",
			requestBody: `{"id": "title2", "type": "gauge"}`,
			wantErr:     false,
		},
		{
			name:        "wrong type",
			requestBody: `{"id": "title1", "type": "testtest"}`,
			wantErr:     true,
		},
		{
			name:        "counter error",
			requestBody: `{"id": "title1", "type": "counter"}`,
			withError:   true,
			wantErr:     true,
		},
		{
			name:        "gauge error",
			requestBody: `{"id": "title2", "type": "gauge"}`,
			withError:   true,
			wantErr:     true,
		},
		{
			name:        "counter missing",
			requestBody: `{"id": "not-found", "type": "counter"}`,
			wantErr:     true,
		},
		{
			name:        "gauge missing",
			requestBody: `{"id": "not-found", "type": "gauge"}`,
			wantErr:     true,
		},
	}
	// Create a mock storage
	mms := &metricstorage.MockErrorStorage{MemStorage: *metricstorage.NewMemStorage()}
	ms := metricstorage.NewSavableModelStorage(mms)
	i64 := int64(1)
	f64 := float64(1.1)
	_ = ms.SaveCommonModels(context.Background(), []apimodels.Metrics{
		{
			ID:    "title1",
			MType: "counter",
			Delta: &i64,
		},
		{
			ID:    "title2",
			MType: "gauge",
			Value: &f64,
		},
	})
	h := CheckMetricExistenceHandler(ms)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mms.WithError = tt.withError
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			rc := io.NopCloser(strings.NewReader(tt.requestBody))
			r := httptest.NewRequest(http.MethodPost, "/update", rc)
			c.Request = r
			h(c)

			if (w.Code != http.StatusOK || len(c.Errors) > 0) != tt.wantErr {
				t.Errorf("CheckMetricExistenceHandler() errors = %v, wantErr %v", len(c.Errors), tt.wantErr)
			}
		})
	}
}
