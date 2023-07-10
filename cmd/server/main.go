package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oriolus/notprometheus/internal/handler"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage/memory"
)

func main() {
	mux := chi.NewRouter()

	storage := memory.NewMemStorage()
	metricServer := server.NewServer(storage)
	updateHandler := handler.NewUpdateHandler(metricServer)

	updatePattern := fmt.Sprintf("/update/{%s}/{%s}/{%s}", handler.URLParamMetricType, handler.URLParamName, handler.URLParamValue)
	mux.Post(updatePattern, updateHandler.ServeHTTP)

	err := http.ListenAndServe("localhost:8084", mux)
	if err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}
