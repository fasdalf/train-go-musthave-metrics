package config

import (
	goflag "flag"
	"fmt"
	env "github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

const (
	defaultAddress = ":8080"
)

type Config struct {
	Addr string `env:"ADDRESS"`
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

	// Env. variables. This should take over the command line. Bad practice as I know.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
