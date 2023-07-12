package main

import (
	"flag"
	"time"
)

type config struct {
	address        string
	base           string
	reportInterval time.Duration
	pollInterval   time.Duration
}

const (
	addressDefault        = "http://localhost:8080"
	baseDefault           = "metrics"
	reportIntervalDefault = 10 * time.Second
	pollIntervalDefault   = 2 * time.Second
)

var defaultConfig = &config{
	address:        addressDefault,
	base:           baseDefault,
	reportInterval: reportIntervalDefault,
	pollInterval:   pollIntervalDefault,
}

func parseFlags() *config {
	cfg := &config{}

	flag.StringVar(&cfg.address, "a", addressDefault, "listening address")
	flag.StringVar(&cfg.base, "b", baseDefault, "base")
	flag.DurationVar(&cfg.reportInterval, "r", reportIntervalDefault, "report interval")
	flag.DurationVar(&cfg.pollInterval, "p", pollIntervalDefault, "poll interval")

	flag.Parse()

	return cfg
}
