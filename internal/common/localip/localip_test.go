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

func TestGetFreePort(t *testing.T) {
	gotPort, err := GetFreePort()
	if err != nil {
		t.Errorf("GetFreePort() error = %v", err)
		return
	}
	if gotPort == 0 {
		t.Errorf("GetFreePort() gotPort = %v", gotPort)
	}
}
