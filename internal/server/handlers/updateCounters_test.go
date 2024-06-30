package handlers

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type args struct {
		s metricstorage.Storage
	}
	type want struct {
		statusCode int
		counters   int
		gauges     int
	}
	tests := []struct {
		name string
		url  string
		args args
		want want
	}{
		{
			name: "empty",
			url:  "/",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusNotFound, counters: 0, gauges: 0},
		},
		{
			name: "gauge",
			url:  "/gauge/some-metric/10.001",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusOK, counters: 0, gauges: 1},
		},
		{
			name: "counter",
			url:  "/counter/some-metric/10",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusOK, counters: 1, gauges: 0},
		},
		{
			name: "NaN",
			url:  "/counter/some-metric/NaN",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusBadRequest, counters: 0, gauges: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			handler := UpdateMetricHandler(tt.args.s)
			handler.ServeHTTP(w, r)
			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.gauges, len(tt.args.s.ListGauges()))
			assert.Equal(t, tt.want.counters, len(tt.args.s.ListCounters()))
		})
	}
}
