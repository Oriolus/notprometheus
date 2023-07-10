package main

import (
	"context"
	"fmt"

	"github.com/oriolus/notprometheus/internal/collector"
)

func main() {
	client, err := collector.NewClient("http://localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	server, _ := collector.NewServer(client)
	err = server.Run(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
