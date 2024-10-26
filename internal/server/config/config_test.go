package config

import "testing"

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()

	if cfg.Addr == "" {
		t.Errorf("Expected ADDRESS to be %s, but got %s", defaultAddress, cfg.Addr)
	}
}
