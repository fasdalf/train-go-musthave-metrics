package config

import (
	goflag "flag"
	"fmt"
	env "github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
	"log/slog"
	"os"
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
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	// pflag handles --help itself.

	// Env. variables. This should take over the command line. Bad practice as I know.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
