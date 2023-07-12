package main

import (
	"context"
	"fmt"

	"github.com/oriolus/notprometheus/internal/collector"
	"github.com/oriolus/notprometheus/internal/collector/sender/http"
)

func main() {
	cfg := parseFlags()
	url := cfg.address + "/" + cfg.base

	client, err := http.NewClient(url)
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
