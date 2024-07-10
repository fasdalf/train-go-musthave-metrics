package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name         string
		vType, vName string
		want         want
	}{
		{
			name:  "empty",
			vType: "", vName: "",
			want: want{statusCode: http.StatusBadRequest, body: fmt.Sprintf(
				"Invalid type, only %s and %s supported\n",
				constants.GaugeStr,
				constants.CounterStr,
			)},
		},
		{
			name:  "valid gauge",
			vType: "gauge", vName: "median",
			want: want{statusCode: http.StatusOK, body: "10.001"},
		},
		{
			name:  "valid counter",
			vType: "counter", vName: "amount",
			want: want{statusCode: http.StatusOK, body: "10"},
		},
		{
			name:  "invalid gauge",
			vType: "gauge", vName: "other",
			want: want{statusCode: http.StatusNotFound, body: "metric \"other\" not found\n"},
		},
		{
			name:  "invalid counter",
			vType: "counter", vName: "other",
			want: want{statusCode: http.StatusNotFound, body: "metric \"other\" not found\n"},
		},
	}
	ms := metricstorage.NewMemStorage()
	ms.UpdateGauge("median", 10.001)
	ms.UpdateCounter("amount", 10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/unused", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{
				{Key: "type", Value: tt.vType},
				{Key: "name", Value: tt.vName},
			}
			handler := NewViewStatsHandler(ms)
			handler(c)
			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.body, w.Body.String(), "Содержание ответа не совпадает с ожидаемым")
		})
	}
}
