package main

import (
	"context"
	"fmt"

	"github.com/oriolus/notprometheus/internal/collector"
	"github.com/oriolus/notprometheus/internal/collector/sender/http"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Printf("error while getting config: %s\r\n", err.Error())
	}

	fmt.Printf("Starting agent with config: %s\r\n", cfg)

	url := "http://" + cfg.address
	if cfg.base != "" {
		url += "/" + cfg.base
	}

	client, err := http.NewJSONClient(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	server, _ := collector.NewServer(client, cfg.pollInterval, cfg.reportInterval)
	err = server.Run(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
