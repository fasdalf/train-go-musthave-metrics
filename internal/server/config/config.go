package config

import (
	"crypto/rsa"
	goflag "flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/configfile"

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
	Addr                     string          `env:"ADDRESS" json:"address"`
	StorageFileName          string          `env:"FILE_STORAGE_PATH" json:"store_file"`
	StorageFileStoreInterval int             `env:"STORE_INTERVAL" json:"store_interval"`
	StorageFileRestore       bool            `env:"RESTORE" json:"restore"`
	StorageDBDSN             string          `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey                  string          `env:"KEY" json:"key"`
	RSAKeyFile               string          `env:"CRYPTO_KEY" json:"crypto_key"`
	RSAKey                   *rsa.PrivateKey `env:"-" json:"-"`
	TrustedSubnetCIDR        string          `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	TrustedSubnet            *net.IPNet      `env:"-" json:"-"`
}

var config *Config = &Config{
	Addr:                     defaultAddress,
	StorageFileStoreInterval: defaultStorageFileStoreInterval,
	StorageFileRestore:       defaultStorageFileRestore,
}

func GetConfig() Config {
	return *config
}

func init() {
	configfile.ParseFile(config)

	// Flags
	flag.StringVarP(&config.Addr, "address", "a", config.Addr, "The address to listen on for HTTP requests")
	flag.StringVarP(&config.StorageFileName, "filestoragepath", "f", config.StorageFileName, "A path and file name to store JSON data. Leave empty to disable file storage.")
	flag.IntVarP(&config.StorageFileStoreInterval, "storeinterval", "i", config.StorageFileStoreInterval, "Interval in seconds to dump metrics to JSON file. \"0\" is on every change.")
	flag.BoolVarP(&config.StorageFileRestore, "restore", "r", config.StorageFileRestore, "Metrics restore from file on startup.")
	flag.StringVarP(&config.StorageDBDSN, "databasedsn", "d", config.StorageDBDSN, "Postgres PGX DSN to use DB storage. Disabled when empty.")
	flag.StringVarP(&config.HashKey, "key", "k", config.HashKey, "Key for signature Hash header.  If not provided, will not sign the request.")
	flag.StringVarP(&config.RSAKeyFile, "cryptokey", "b", config.RSAKeyFile, "Private key for body decryption. If not provided, will not decrypt the request.")
	flag.StringVarP(&config.TrustedSubnetCIDR, "trustedsubnet", "t", config.TrustedSubnetCIDR, "Trusted agents subnet CIDR. If not provided all requests are accepted.")
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
	if config.TrustedSubnetCIDR != "" {
		var err error
		if _, config.TrustedSubnet, err = net.ParseCIDR(config.TrustedSubnetCIDR); err != nil {
			panic(err)
		}
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
