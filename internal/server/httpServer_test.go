package server

import (
	"bytes"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fileStorageMock struct{}

func (*fileStorageMock) SaveWithInterval() error {
	return nil
}

func TestUpdateMetricIntegrational(t *testing.T) {
	type args struct {
		s handlers.Storage
	}
	type want struct {
		statusCode int
		counters   int
		gauges     int
	}
	tests := []struct {
		name string
		url  string
		body []byte
		args args
		want want
	}{
		{
			name: "empty",
			url:  "/update/",
			body: []byte("body"),
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusBadRequest, counters: 0, gauges: 0},
		},
		{
			name: "old gauge",
			url:  "/update/gauge/some-metric/10.001",
			body: []byte{},
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusOK, counters: 0, gauges: 1},
		},
		{
			name: "old counter",
			url:  "/update/counter/some-metric/10",
			body: []byte{},
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusOK, counters: 1, gauges: 0},
		},
		{
			name: "old NaN",
			url:  "/update/counter/some-metric/NaN",
			body: []byte{},
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusBadRequest, counters: 0, gauges: 0},
		},
		{
			name: "gauge",
			url:  "/update/",
			body: []byte(`{"id": "some-metric","type": "gauge", "value": 10.001}`),
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusOK, counters: 0, gauges: 1},
		},
		{
			name: "counter",
			url:  "/update/",
			body: []byte(`{"id": "some-metric","type": "counter", "delta": 10}`),
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusOK, counters: 1, gauges: 0},
		},
		{
			name: "NaN",
			url:  "/update/",
			body: []byte(`{"id": "some-metric","type": "gauge", "delta": "NaN"}`),
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusBadRequest, counters: 0, gauges: 0},
		},
		{
			name: "batch success",
			url:  "/updates/",
			body: []byte(`[{"id": "some-metric","type": "counter", "delta": 10},{"id": "some-gauge","type": "gauge", "value": 100.01}]`),
			args: args{metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())},
			want: want{statusCode: http.StatusOK, counters: 1, gauges: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRoutingEngine(tt.args.s, nil, nil, retryattempt.NewOneAttemptRetryer())

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer(tt.body))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.gauges, len(tt.args.s.ListGauges()))
			assert.Equal(t, tt.want.counters, len(tt.args.s.ListCounters()))
		})
	}
}

func TestViewMetricIntegrational(t *testing.T) {
	type want struct {
		statusCode int
		json       string
		plain      string
	}
	tests := []struct {
		name   string
		method string
		url    string
		body   []byte
		want   want
	}{
		{
			name:   "not found GET",
			method: http.MethodGet,
			url:    "/value/",
			body:   []byte("body"),
			want:   want{statusCode: http.StatusNotFound},
		},
		{
			name:   "old gauge",
			method: http.MethodGet,
			url:    "/value/gauge/Floating",
			body:   []byte{},
			want:   want{statusCode: http.StatusOK, plain: "100.001"},
		},
		{
			name:   "old counter",
			method: http.MethodGet,
			url:    "/value/counter/Integral",
			body:   []byte{},
			want:   want{statusCode: http.StatusOK, plain: "10"},
		},
		{
			name:   "old invalid type",
			method: http.MethodGet,
			url:    "/value/missing/ignored",
			body:   []byte{},
			want:   want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "old gauge not found",
			method: http.MethodGet,
			url:    "/value/gauge/notFloating",
			body:   []byte{},
			want:   want{statusCode: http.StatusNotFound},
		},
		{
			name:   "old counter not found",
			method: http.MethodGet,
			url:    "/view/counter/notIntegral",
			body:   []byte{},
			want:   want{statusCode: http.StatusNotFound},
		},

		{
			name:   "invalid POST body",
			method: http.MethodPost,
			url:    "/value",
			body:   []byte("body"),
			want:   want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "gauge v2",
			method: http.MethodPost,
			url:    "/value",
			body:   []byte(`{"type":"gauge","id":"Floating"}`),
			want:   want{statusCode: http.StatusOK, json: `{"type":"gauge","id":"Floating","value":100.001}`},
		},
		{
			name:   "counter v2",
			method: http.MethodPost,
			url:    "/value",
			body:   []byte(`{"type":"counter","id":"Integral"}`),
			want:   want{statusCode: http.StatusOK, json: `{"type":"counter","id":"Integral","delta":10}`},
		},
		{
			name:   "old invalid type",
			method: http.MethodPost,
			url:    "/value",
			body:   []byte(`{"type":"invalid","id":"useless"}`),
			want:   want{statusCode: http.StatusBadRequest},
		},
		{
			name:   "old gauge not found",
			method: http.MethodPost,
			url:    "/value",
			body:   []byte(`{"type":"gauge","id":"notFloating"}`),
			want:   want{statusCode: http.StatusNotFound},
		},
		{
			name:   "old counter not found",
			method: http.MethodPost,
			url:    "/value",
			body:   []byte(`{"type":"counter","id":"notIntegral"}`),
			want:   want{statusCode: http.StatusNotFound},
		},
	}

	ms := metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())
	ms.UpdateGauge("Floating", 100.001)
	ms.UpdateCounter("Integral", 10)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRoutingEngine(ms, nil, nil, retryattempt.NewOneAttemptRetryer())

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.body))
			router.ServeHTTP(w, req)

			if tt.want.json != "" {
				assert.JSONEq(t, tt.want.json, w.Body.String())
			}
			if tt.want.json != "" {
				assert.JSONEq(t, tt.want.json, w.Body.String())
			}

			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}
