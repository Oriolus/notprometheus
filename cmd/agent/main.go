package main

import (
	"context"
	"github.com/oriolus/notprometheus/internal/collector"
	"github.com/oriolus/notprometheus/internal/collector/sender/http"
	"github.com/oriolus/notprometheus/internal/logger"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		logger.Fatalf("error while getting config: %s\r\n", err.Error())
	}

	logger.Infof("Starting agent with config: %s\r\n", cfg)

	url := "http://" + cfg.address
	if cfg.base != "" {
		url += "/" + cfg.base
	}

	client, err := http.NewJSONClient(url)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	server, _ := collector.NewServer(client, cfg.pollInterval, cfg.reportInterval)
	err = server.Run(context.Background())
	if err != nil {
		logger.Fatalf(err.Error())
	}
}
