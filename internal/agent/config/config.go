package config

import (
	goflag "flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
	"log/slog"
	"os"
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
}

var config *Config

func GetConfig() Config {
	return *config
}

func init() {
	config = &Config{}
	// Flags
	flag.StringVarP(&config.Addr, "address", "a", defaultAddress, fmt.Sprintf("The address to listen on for HTTP requests. Default is \"%s\"", defaultAddress))
	flag.IntVarP(&config.PollInterval, "pollinterval", "p", defaultPollInterval, fmt.Sprintf("Metrics poll interval in seconds. Default is \"%d\"", defaultPollInterval))
	flag.IntVarP(&config.ReportInterval, "reportInterval", "r", defaultPollInterval, fmt.Sprintf("Report to server interval in seconds. Default is \"%d\"", defaultReportInterval))
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	// Env. variables. This should take over the command line. Bad practice as I know.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
