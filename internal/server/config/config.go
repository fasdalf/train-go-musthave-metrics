package config

import (
	goflag "flag"
	"fmt"
	flag "github.com/spf13/pflag"
)

const (
	defaultAddress = ":8080"
)

type Config struct {
	Addr string
}

var config *Config

func GetConfig() Config {
	return *config
}
func init() {
	config = &Config{}
	// Flags
	flag.StringVarP(&config.Addr, "address", "a", defaultAddress, fmt.Sprintf("The address to listen on for HTTP requests. Default is \"%s\"", defaultAddress))
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
