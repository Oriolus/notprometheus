package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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

	addressEnvParamName        = "SERVER_ADDRESS"
	baseEnvParamName           = "BASE_URL"
	pollIntervalEnvParamName   = "POLL_INTERVAL"
	reportIntervalEnvParamName = "REPORT_INTERVAL"
)

var defaultConfig = &config{
	address:        addressDefault,
	base:           baseDefault,
	reportInterval: reportIntervalDefault,
	pollInterval:   pollIntervalDefault,
}

func parseFlags() (*config, error) {
	cfg := &config{}

	flag.StringVar(&cfg.address, "a", addressDefault, "listening address")
	flag.StringVar(&cfg.base, "b", baseDefault, "base")
	flag.DurationVar(&cfg.reportInterval, "r", reportIntervalDefault, "report interval")
	flag.DurationVar(&cfg.pollInterval, "p", pollIntervalDefault, "poll interval")

	flag.Parse()

	if val, ok := os.LookupEnv(addressEnvParamName); ok {
		cfg.address = val
	}

	if val, ok := os.LookupEnv(baseEnvParamName); ok {
		cfg.base = val
	}

	if val, ok := os.LookupEnv(pollIntervalEnvParamName); ok {
		pollInterval, err := getDuration(val, pollIntervalEnvParamName)
		if err != nil {
			return nil, err
		}
		cfg.pollInterval = pollInterval
	}

	if val, ok := os.LookupEnv(reportIntervalEnvParamName); ok {
		reportInterval, err := getDuration(val, reportIntervalEnvParamName)
		if err != nil {
			return nil, err
		}
		cfg.reportInterval = reportInterval
	}

	return cfg, nil
}

func getDuration(val, paramName string) (time.Duration, error) {
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0 * time.Second, fmt.Errorf("erorr while parsing param %s. Got value: %s. Error: %s", paramName, val, err.Error())
	}

	return time.Duration(intVal) * time.Second, nil
}
