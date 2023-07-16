package main

import (
	"flag"
	"fmt"
	"os"
)

type config struct {
	address string
	base    string
}

func (c *config) String() string {
	return fmt.Sprintf("address: %s, base: %s", c.address, c.base)
}

const (
	addressDefault = "localhost:8080"
	baseDefault    = ""

	addressEnvParamName = "SERVER_ADDRESS"
	baseEnvParamName    = "BASE_URL"
)

var defaultConfig = &config{
	address: addressDefault,
	base:    baseDefault,
}

func parseFlags() *config {
	cfg := &config{}

	flag.StringVar(&cfg.address, "a", addressDefault, "listening address")
	flag.StringVar(&cfg.base, "b", baseDefault, "base")

	flag.Parse()

	if val, ok := os.LookupEnv(addressEnvParamName); ok {
		cfg.address = val
	}

	if val, ok := os.LookupEnv(baseEnvParamName); ok {
		cfg.base = val
	}

	return cfg
}
