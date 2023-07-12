package main

import "flag"

type config struct {
	address string
	base    string
}

const (
	addressDefault = "localhost:8080"
	baseDefault    = "metrics"
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

	return cfg
}
