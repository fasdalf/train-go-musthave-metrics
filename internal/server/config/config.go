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
	defaultStorageFileName          = "data.json"
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
	flag.StringVarP(&config.Addr, "address", "a", defaultAddress, fmt.Sprintf("The address to listen on for HTTP requests. Default is \"%s\"", defaultAddress))
	flag.StringVarP(&config.StorageFileName, "filestoragepath", "f", defaultStorageFileName, fmt.Sprintf("A path and file name to store data. Default is \"%s\"", defaultStorageFileName))
	flag.IntVarP(&config.StorageFileStoreInterval, "storeinterval", "i", defaultStorageFileStoreInterval, fmt.Sprintf("Metrics store interval in seconds. \"0\" is on every change. Default is \"%d\"", defaultStorageFileStoreInterval))
	flag.BoolVarP(&config.StorageFileRestore, "restore", "r", defaultStorageFileRestore, fmt.Sprintf("Metrics restore from file on startup. Default is \"%s\"", defaultStorageFileRestore))
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	// Env. variables. This should take over the command line. Bad practice as I know.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
