package localip

import (
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	if len(ip) == 0 {
		t.Errorf("Expected a non-empty IP address, got %v", ip)
	}
}
