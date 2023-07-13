package main

import (
	"fmt"
	"net/http"

	"github.com/oriolus/notprometheus/internal/handler"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage/memory"
)

func main() {
	storage := memory.NewMemStorage()
	metricServer := server.NewServer(storage)
	updateHandler := handler.NewUpdateHandler(metricServer)

	serveMux := http.NewServeMux()
	serveMux.Handle("/update/", updateHandler)

	err := http.ListenAndServe("localhost:8080", serveMux)
	if err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}
