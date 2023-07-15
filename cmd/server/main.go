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

func ChiRouter(cfg *config) chi.Router {
	mux := chi.NewRouter()
	mux.Use(middleware.DefaultLogger)

	storage := memory.NewMemStorage()
	metricServer := server.NewServer(storage)
	updateHandler, _ := handler.NewUpdateHandler(metricServer)
	getAllHandler, _ := handler.NewGetAllHandler(metricServer)
	getMetricValue, _ := handler.NewGetMetricHandler(metricServer)

	updatePattern := fmt.Sprintf("/update/{%s}/{%s}/{%s}", handler.URLParamMetricType, handler.URLParamName, handler.URLParamValue)
	updatePattern = addBase(cfg, updatePattern)
	mux.Post(updatePattern, updateHandler.ServeHTTP)

	getAllPattern := addBase(cfg, "/")
	mux.Get(getAllPattern, getAllHandler.ServeHTTP)

	getMetricPattern := fmt.Sprintf("/value/{%s}/{%s}", handler.URLParamMetricType, handler.URLParamName)
	getMetricPattern = addBase(cfg, getMetricPattern)
	mux.Get(getMetricPattern, getMetricValue.ServeHTTP)
	return mux
}

func main() {
	cfg := parseFlags()
	mux := ChiRouter(cfg)

	fmt.Printf("Listening address %s with base %s", cfg.address, cfg.base)

	err := http.ListenAndServe(cfg.address, mux)
	if err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}

func addBase(cfg *config, url string) string {
	if cfg == nil || cfg.base == "" {
		return url
	}
	return "/" + cfg.base + url
}
