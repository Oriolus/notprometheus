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

func (c *config) String() string {
	return fmt.Sprintf("address: %s, base: %s, report interval: %s, poll interval: %s",
		c.address, c.base, c.reportInterval.String(), c.pollInterval.String())
}

const (
	addressDefault           = "localhost:8080"
	baseDefault              = ""
	reportIntervalSecDefault = 10
	pollIntervalSecDefault   = 2

	addressEnvParamName        = "SERVER_ADDRESS"
	baseEnvParamName           = "BASE_URL"
	pollIntervalEnvParamName   = "POLL_INTERVAL"
	reportIntervalEnvParamName = "REPORT_INTERVAL"
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
