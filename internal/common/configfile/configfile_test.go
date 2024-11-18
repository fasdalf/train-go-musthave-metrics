// Package configfile - config file package
package configfile

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	oldEnv := os.Getenv(configEnv)
	defer os.Setenv(configEnv, oldEnv)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	c := struct {
		Address string
	}{}
	const envFile = "env.file.json"

	os.Setenv(configEnv, "")
	os.Args = []string{"executable"}

	ParseFile(&c)

	if c.Address != "" {
		t.Errorf("c.Address should be empty, got \"%v\"", c.Address)
	}

	os.Setenv(configEnv, envFile)

	os.WriteFile(envFile, []byte(`{"address":":1080"}`), 0644)
	defer os.Remove(envFile)

	ParseFile(&c)

	if c.Address != ":1080" {
		t.Errorf("c.Address should be \":1080\", got \"%v\"", c.Address)
	}
}

func TestGetConfigFileName(t *testing.T) {
	oldEnv := os.Getenv(configEnv)
	defer os.Setenv(configEnv, oldEnv)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	const envFile = "env.file.json"
	const argFile = "arg.file.json"

	os.Setenv(configEnv, envFile)
	os.Args = []string{}

	r := getConfigFileName()

	if r != envFile {
		t.Errorf("getConfigFileName() returns %s, want %s", r, envFile)
	}

	os.Setenv(configEnv, "")
	os.Args = []string{"executable", fmt.Sprintf("-%s=%s", configFlag, argFile)}

	r = getConfigFileName()

	if r != argFile {
		t.Errorf("getConfigFileName() returns %s, want %s", r, argFile)
	}
}
