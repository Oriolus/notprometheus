package main

import (
	"flag"
	"fmt"
	"time"
)

type config struct {
	address        string
	base           string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func (c *config) String() string {
	return fmt.Sprintf("address: %s, base: %s, report interval: %s, poll interval: %s",
		c.address, c.base, c.reportInterval.String(), c.pollInterval.String())
}

const (
	addressDefault           = "localhost:8080"
	baseDefault              = ""
	reportIntervalSecDefault = 10
	pollIntervalSecDefault   = 2
)

var defaultConfig = &config{
	address:        addressDefault,
	base:           baseDefault,
	reportInterval: reportIntervalSecDefault * time.Second,
	pollInterval:   pollIntervalSecDefault * time.Second,
}

func parseFlags() (*config, error) {
	cfg := &config{}

	var (
		reportIntervalSec int64
		pollIntervalSec   int64
	)

	flag.StringVar(&cfg.address, "a", addressDefault, "listening address")
	flag.StringVar(&cfg.base, "b", baseDefault, "base")
	flag.Int64Var(&reportIntervalSec, "r", reportIntervalSecDefault, "report interval")
	flag.Int64Var(&pollIntervalSec, "p", pollIntervalSecDefault, "poll interval")

	flag.Parse()

	cfg.reportInterval = time.Duration(reportIntervalSec) * time.Second
	cfg.pollInterval = time.Duration(pollIntervalSec) * time.Second

	return cfg, nil
}
