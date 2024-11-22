package handlers

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/gin-gonic/gin"
)

// NewValidateIPHandler validate IP header
// Trust such headers only then they're set by trusted reverse proxy.
func NewValidateIPHandler(tr *net.IPNet) func(c *gin.Context) {
	validateIPHandler := func(c *gin.Context) {
		if tr != nil {
			ipString := c.GetHeader(constants.HeaderRealIP)
			if err := validateIPStringInSubnet(ipString, tr); err != nil {
				slog.Error("header value is invalid", "header", constants.HeaderRealIP, "error", err)
				_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("%s header value is invalid: %s", constants.HeaderRealIP, ipString))
				return
			}

		}

		c.Next()
	}
	return validateIPHandler
}

func validateIPStringInSubnet(addr string, subnet *net.IPNet) error {
	ip := net.ParseIP(addr)
	if ip == nil {
		return fmt.Errorf("\"%s\" is not a valid IP address", addr)
	}
	if subnet == nil {
		return fmt.Errorf("empty subnet")
	}
	if !subnet.Contains(ip) {
		return fmt.Errorf("IP address \"%s\" is not in subnet \"%s\"", addr, subnet.String())
	}
	return nil
}
