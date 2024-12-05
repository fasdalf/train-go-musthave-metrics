package handlers

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/localip"
)

// NewValidateIPHandler validate IP header
// Trust such headers only then they're set by trusted reverse proxy.
func NewValidateIPHandler(tr *net.IPNet) func(c *gin.Context) {
	validateIPHandler := func(c *gin.Context) {
		if tr != nil {
			ipString := c.GetHeader(constants.HeaderRealIP)
			if err := localip.ValidateIPStringInSubnet(ipString, tr); err != nil {
				slog.Error("header value is invalid", "header", constants.HeaderRealIP, "error", err)
				_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("%s header value is invalid: %s", constants.HeaderRealIP, ipString))
				return
			}
		}

		c.Next()
	}
	return validateIPHandler
}
