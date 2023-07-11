package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oriolus/notprometheus/internal/handler"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage/memory"
)

func ChiRouter() chi.Router {
	mux := chi.NewRouter()
	mux.Use(middleware.DefaultLogger)

	storage := memory.NewMemStorage()
	metricServer := server.NewServer(storage)
	updateHandler, _ := handler.NewUpdateHandler(metricServer)
	getAllHandler, _ := handler.NewGetAllHandler(metricServer)
	getMetricValue, _ := handler.NewGetMetricHandler(metricServer)

	updatePattern := fmt.Sprintf("/update/{%s}/{%s}/{%s}", handler.URLParamMetricType, handler.URLParamName, handler.URLParamValue)
	mux.Post(updatePattern, updateHandler.ServeHTTP)

	mux.Get("/", getAllHandler.ServeHTTP)

	getMetricPattern := fmt.Sprintf("/value/{%s}/{%s}", handler.URLParamMetricType, handler.URLParamName)
	mux.Get(getMetricPattern, getMetricValue.ServeHTTP)
	return mux
}

func main() {
	mux := ChiRouter()
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}
