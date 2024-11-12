package config

import (
	"crypto/rsa"
	goflag "flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"

	env "github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

const (
	defaultAddress                  = ":8080"
	defaultStorageFileStoreInterval = 300
	defaultStorageFileRestore       = true
)

type Config struct {
	Addr                     string `env:"ADDRESS"`
	StorageFileName          string `env:"FILE_STORAGE_PATH"`
	StorageFileStoreInterval int    `env:"STORE_INTERVAL"`
	StorageFileRestore       bool   `env:"RESTORE"`
	StorageDBDSN             string `env:"DATABASE_DSN"`
	HashKey                  string `env:"KEY"`
	RSAKeyFile               string `env:"CRYPTO_KEY"`
	RSAKey                   *rsa.PrivateKey
}

var config *Config

func GetConfig() Config {
	return *config
}

func init() {
	config = &Config{}
	// Flags
	flag.StringVarP(&config.Addr, "address", "a", defaultAddress, "The address to listen on for HTTP requests")
	flag.StringVarP(&config.StorageFileName, "filestoragepath", "f", "", "A path and file name to store JSON data. Leave empty to disable file storage.")
	flag.IntVarP(&config.StorageFileStoreInterval, "storeinterval", "i", defaultStorageFileStoreInterval, "Interval in seconds to dump metrics to JSON file. \"0\" is on every change.")
	flag.BoolVarP(&config.StorageFileRestore, "restore", "r", defaultStorageFileRestore, "Metrics restore from file on startup.")
	flag.StringVarP(&config.StorageDBDSN, "databasedsn", "d", "", "Postgres PGX DSN to use DB storage. Disabled when empty.")
	flag.StringVarP(&config.HashKey, "key", "k", "", "Key for signature Hash header.  If not provided, will not sign the request.")
	flag.StringVarP(&config.RSAKeyFile, "crypto-key", "c", "", "Private key for body decryption. If not provided, will not decrypt the request.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	// pflag handles --help itself.

	// Env. variables. This should take over the command line. Bad practice as I know.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}
	if config.RSAKeyFile != "" {
		var err error
		if config.RSAKey, err = rsacrypt.FileToPrivateKey(config.RSAKeyFile); err != nil {
			panic(err)
		}
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
