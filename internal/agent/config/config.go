package config

import (
	"crypto/rsa"
	goflag "flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/configfile"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

const (
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	Addr           string         `env:"ADDRESS" json:"address"`
	Protocol       string         `env:"PROTOCOL" json:"protocol"`
	PollInterval   int            `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval int            `env:"REPORT_INTERVAL" json:"report_interval"`
	HashKey        string         `env:"KEY" json:"key"`
	RateLimit      int            `env:"RATE_LIMIT" json:"rate_limit"`
	RSAKeyFile     string         `env:"CRYPTO_KEY" json:"crypto_key"`
	RSAKey         *rsa.PublicKey `env:"-" json:"-"`
}

var config *Config = &Config{
	Addr:           defaultAddress,
	PollInterval:   defaultPollInterval,
	ReportInterval: defaultReportInterval,
	HashKey:        "",
	RateLimit:      1,
	RSAKeyFile:     "",
	RSAKey:         nil,
}

func GetConfig() Config {
	return *config
}

func init() {
	configfile.ParseFile(config)

	// Flags
	flag.StringVarP(&config.Addr, "address", "a", config.Addr, "The address to listen on for HTTP requests.")
	flag.StringVarP(&config.Protocol, "protocol", "g", config.Protocol, "Protocol to use instead of HTTP. Now supported: grpc.")
	flag.IntVarP(&config.PollInterval, "pollinterval", "p", config.PollInterval, "Metrics poll interval in seconds.")
	flag.IntVarP(&config.ReportInterval, "reportInterval", "r", config.ReportInterval, "Report to server interval in seconds.")
	flag.StringVarP(&config.HashKey, "key", "k", config.HashKey, "Key for signature Hash header.  If not provided, will not sign the request.")
	flag.IntVarP(&config.RateLimit, "ratelimit", "l", config.RateLimit, "Limit number of outcoming requests.")
	flag.StringVarP(&config.RSAKeyFile, "cryptokey", "b", config.RSAKeyFile, "Public key for body encryption. If not provided, will not encrypt the request.")
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
	if config.RSAKeyFile != "" {
		var err error
		if config.RSAKey, err = rsacrypt.FileToPublicKey(config.RSAKeyFile); err != nil {
			panic(err)
		}
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
