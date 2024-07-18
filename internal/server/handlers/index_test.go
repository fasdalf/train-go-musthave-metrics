package handlers

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewIndexHandler(t *testing.T) {
	type args struct {
		s Storage
	}
	tests := []struct {
		name     string
		args     args
		wantBody string
	}{
		{
			name: "empty",
			args: args{metricstorage.NewMemStorage()},
			wantBody: `<html><body>
<table>
<tr><td colspan=2>Gauges</td></tr>
<tr><td colspan=2>counters</td></tr>
</table>
</body></html>
`,
		},
		{
			name: "valid gauge",
			args: args{(func() Storage {
				ms := metricstorage.NewMemStorage()
				ms.UpdateGauge("title1", 20.1)
				return ms
			})()},
			wantBody: `<html><body>
<table>
<tr><td colspan=2>Gauges</td></tr>
<tr><td>title1</td><td>20.1</td></tr>
<tr><td colspan=2>counters</td></tr>
</table>
</body></html>
`,
		},
		{
			name: "valid counter",
			args: args{(func() Storage {
				ms := metricstorage.NewMemStorage()
				ms.UpdateCounter("title2", 201)
				return ms
			})()},
			wantBody: `<html><body>
<table>
<tr><td colspan=2>Gauges</td></tr>
<tr><td colspan=2>counters</td></tr>
<tr><td>title2</td><td>201</td></tr>
</table>
</body></html>
`,
		},
		{
			name: "gauge and counter",
			args: args{(func() Storage {
				ms := metricstorage.NewMemStorage()
				ms.UpdateGauge("title1", 20.1)
				ms.UpdateCounter("title2", 201)
				return ms
			})()},
			wantBody: `<html><body>
<table>
<tr><td colspan=2>Gauges</td></tr>
<tr><td>title1</td><td>20.1</td></tr>
<tr><td colspan=2>counters</td></tr>
<tr><td>title2</td><td>201</td></tr>
</table>
</body></html>
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/unused", nil)
			w := httptest.NewRecorder()
			handler := NewIndexHandler(tt.args.s)
			handler.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, "text/html", w.Header().Get("Content-Type"), "Тип ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.wantBody, w.Body.String(), "Содержание ответа не совпадает с ожидаемым")
		})
	}
}
