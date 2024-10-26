package config

import (
	goflag "flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

const (
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

var config *Config

func GetConfig() Config {
	return *config
}

func init() {
	config = &Config{}
	// Flags
	flag.StringVarP(&config.Addr, "address", "a", defaultAddress, "The address to listen on for HTTP requests.")
	flag.IntVarP(&config.PollInterval, "pollinterval", "p", defaultPollInterval, "Metrics poll interval in seconds.")
	flag.IntVarP(&config.ReportInterval, "reportInterval", "r", defaultPollInterval, "Report to server interval in seconds.")
	flag.StringVarP(&config.HashKey, "key", "k", "", "Key for signature Hash header.  If not provided, will not sign the request.")
	flag.IntVarP(&config.RateLimit, "ratelimit", "l", 1, "Limit number of outcoming requests.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	// pflag handles --help itself.

	// Env. variables. This should take over the command line. Bad practice as I know.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if config.RateLimit < 1 {
		panic("ratelimit must be greater than zero")
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
