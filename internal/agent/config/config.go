package config

import (
	goflag "flag"
	"fmt"
	flag "github.com/spf13/pflag"
)

const (
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	Addr                         string
	PollInterval, ReportInterval int
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
}
