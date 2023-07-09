package main

import (
	"fmt"
	"github.com/oriolus/notprometheus/internal/handler"
	"net/http"

	"github.com/oriolus/notprometheus/internal/metric/server"
	"github.com/oriolus/notprometheus/internal/metric/storage/memory"
)

func main() {
	storage := memory.NewMemStorage()
	metricServer := server.NewServer(storage)
	updateHandler := handler.NewUpdateHandler(metricServer)

	serveMux := http.NewServeMux()
	serveMux.Handle("/update/", updateHandler)

	err := http.ListenAndServe("localhost:8084", serveMux)
	if err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}
