package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMetricValueResponseHandler(t *testing.T) {
	tests := []struct {
		name           string
		metric         any
		withError      bool
		wantStatusCode int
	}{
		{
			name:           "empty context",
			metric:         nil,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "not a metric",
			metric:         "nil",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid type",
			metric: &apimodels.Metrics{
				ID:    "any",
				MType: "any",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "valid counter",
			metric: &apimodels.Metrics{
				ID:    "title1",
				MType: "counter",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "valid gauge",
			metric: &apimodels.Metrics{
				ID:    "title1",
				MType: "gauge",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid counter",
			metric: &apimodels.Metrics{
				ID:    "title",
				MType: "counter",
			},
			withError:      true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "invalid gauge",
			metric: &apimodels.Metrics{
				ID:    "title",
				MType: "gauge",
			},
			withError:      true,
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	// Create a mock storage
	mms := &metricstorage.MockErrorStorage{MemStorage: *metricstorage.NewMemStorage()}
	ms := metricstorage.NewSavableModelStorage(mms)
	ms.UpdateGauge("title1", 20.1)
	ms.UpdateCounter("title2", 201)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mms.WithError = tt.withError
			r := httptest.NewRequest(http.MethodPost, "/unused", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			if tt.metric != nil {
				c.Set(contextMetricResponseKey, tt.metric)
			}
			handler := MetricValueResponseHandler(ms)
			handler(c)
			assert.Equal(t, tt.wantStatusCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}
