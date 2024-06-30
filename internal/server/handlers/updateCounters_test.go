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
		args args
		want want
	}{
		{
			name: "single one",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusNotFound, counters: 0, gauges: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()
			handler := UpdateMetricHandler(tt.args.s)
			handler.ServeHTTP(w, r)
			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.gauges, len(tt.args.s.ListGauges()))
			assert.Equal(t, tt.want.counters, len(tt.args.s.ListCounters()))
		})
	}
}
