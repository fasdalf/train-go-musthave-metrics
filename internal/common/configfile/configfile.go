// Package configfile - config file package
package configfile

import (
	"encoding/json"
	"os"

	flag "github.com/spf13/pflag"
)

const (
	configEnv      = "CONFIG"
	configFlag     = "c"
	configFlagFull = "config"
)

func init() {
	var s string
	flag.StringVarP(&s, configFlagFull, configFlag, "", "Path to config file. Not required.")
}

func ParseFile(config any) {
	fn := getConfigFileName()
	if fn == "" {
		return
	}

	data, _ := os.ReadFile(fn)
	json.Unmarshal(data, &config)
}

func getConfigFileName() (f string) {
	f = os.Getenv(configEnv)
	if f != "" {
		return
	}

	cl := flag.NewFlagSet(configFlagFull, flag.ContinueOnError)
	cl.StringVarP(&f, configFlagFull, configFlag, f, "")
	cl.Parse(os.Args[1:])
	return
}
