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

	base := ""
	if cfg.base != "" {
		base = "/" + cfg.base
	}
	updatePattern := fmt.Sprintf("%s/update/{%s}/{%s}/{%s}", base, handler.URLParamMetricType, handler.URLParamName, handler.URLParamValue)
	mux.Post(updatePattern, updateHandler.ServeHTTP)

	mux.Get(fmt.Sprintf("%s/", cfg.base), getAllHandler.ServeHTTP)

	getMetricPattern := fmt.Sprintf("%s/value/{%s}/{%s}", base, handler.URLParamMetricType, handler.URLParamName)
	mux.Get(getMetricPattern, getMetricValue.ServeHTTP)
	return mux
}

func main() {
	cfg := parseFlags()
	mux := ChiRouter(cfg)

	fmt.Printf("Starting server with config: %s\r\n", cfg)

	err := http.ListenAndServe(cfg.address, mux)
	if err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}
