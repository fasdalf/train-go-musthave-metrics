package localip

import (
	"net"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	if len(ip) == 0 {
		t.Errorf("Expected a non-empty IP address, got %v", ip)
	}
}

func TestValidateIPStringInSubnet(t *testing.T) {
	validIP, ValidSubnet, _ := net.ParseCIDR("192.168.1.10/24")

	// Test cases
	tests := []struct {
		name      string
		addr      string
		subnet    *net.IPNet
		wantError bool
	}{
		{
			name:      "ok",
			addr:      validIP.String(),
			subnet:    ValidSubnet,
			wantError: false,
		},
		{
			name:      "ok",
			addr:      "",
			subnet:    ValidSubnet,
			wantError: true,
		},
		{
			name:      "ok",
			addr:      validIP.String(),
			subnet:    nil,
			wantError: true,
		},
		{
			name:      "ok",
			addr:      "192.168.2.1",
			subnet:    ValidSubnet,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateIPStringInSubnet(tt.addr, tt.subnet); (err != nil) != tt.wantError {
				t.Errorf("ValidateIPStringInSubnet() error = %v, wantErr %v", err, tt.wantError)
			}
		})
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
