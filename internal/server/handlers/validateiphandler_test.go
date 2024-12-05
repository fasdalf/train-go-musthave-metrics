package handlers

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateIPHandler(t *testing.T) {
	validIP, ValidSubnet, _ := net.ParseCIDR("192.168.1.10/24")

	tests := []struct {
		name           string
		addr           string
		wantStatusCode int
	}{
		{
			name:           "ok",
			addr:           validIP.String(),
			wantStatusCode: 200,
		},
		{
			name:           "err",
			addr:           "nil",
			wantStatusCode: 403,
		},
	}

	h := NewValidateIPHandler(ValidSubnet)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/unused", nil)
			r.Header.Add(constants.HeaderRealIP, tt.addr)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			h(c)
			assert.Equal(t, tt.wantStatusCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}
