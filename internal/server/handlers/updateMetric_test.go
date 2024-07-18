package handlers

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type args struct {
		s Storage
	}
	type want struct {
		statusCode int
		counters   int
		gauges     int
	}
	tests := []struct {
		name                 string
		vType, vName, vValue string
		args                 args
		want                 want
	}{
		{
			name:  "empty",
			vType: "", vName: "", vValue: "",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusNotFound, counters: 0, gauges: 0},
		},
		{
			name:  "valid gauge",
			vType: "gauge", vName: "some-metric", vValue: "10.001",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusOK, counters: 0, gauges: 1},
		},
		{
			name:  "valid counter",
			vType: "counter", vName: "some-metric", vValue: "10",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusOK, counters: 1, gauges: 0},
		},
		{
			name:  "NaN counter",
			vType: "counter", vName: "some-metric", vValue: "NaN",
			args: args{metricstorage.NewMemStorage()},
			want: want{statusCode: http.StatusBadRequest, counters: 0, gauges: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/unused", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{
				{Key: "type", Value: tt.vType},
				{Key: "name", Value: tt.vName},
				{Key: "value", Value: tt.vValue},
			}
			handler := NewUpdateMetricHandler(tt.args.s)
			handler(c)
			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.gauges, len(tt.args.s.ListGauges()), "Got unexpected amount of gauges")
			assert.Equal(t, tt.want.counters, len(tt.args.s.ListCounters()), "Got unexpected amount of counters")
		})
	}
}
