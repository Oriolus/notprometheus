package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oriolus/notprometheus/cmd/server/middleware"
	"github.com/oriolus/notprometheus/internal/handler"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage/memory"
)

func ChiRouter(cfg *config) chi.Router {
	mux := chi.NewRouter()
	mux.Use(middleware.WithLogging)

	storage := memory.NewMemStorage()
	metricServer := server.NewServer(storage)

	base := ""
	if cfg.base != "" {
		base = "/" + cfg.base
	}
	updateHandler, _ := handler.NewUpdateHandler(metricServer)
	updatePattern := fmt.Sprintf("%s/update/{%s}/{%s}/{%s}", base, handler.URLParamMetricType, handler.URLParamName, handler.URLParamValue)
	mux.Post(updatePattern, updateHandler.ServeHTTP)

	getAllHandler, _ := handler.NewGetAllHandler(metricServer)
	mux.Get(fmt.Sprintf("%s/", cfg.base), getAllHandler.ServeHTTP)

	getMetricPattern := fmt.Sprintf("%s/value/{%s}/{%s}", base, handler.URLParamMetricType, handler.URLParamName)
	getMetricValue, _ := handler.NewGetMetricHandler(metricServer)
	mux.Get(getMetricPattern, getMetricValue.ServeHTTP)

	updateJSONHandler, _ := handler.NewUpdateJSONedHandler(metricServer)
	updateJSONPattern := fmt.Sprintf("%s/update", base)
	mux.Post(updateJSONPattern, updateJSONHandler.ServeHTTP)

	getJSONHandler, _ := handler.NewGetJSONedHandler(metricServer)
	getJSONPattern := fmt.Sprintf("%s/value", base)
	mux.Get(getJSONPattern, getJSONHandler.ServeHTTP)

	return mux
}

func main() {
	cfg := parseFlags()
	mux := ChiRouter(cfg)

	fmt.Printf("Starting server with config: %s\r\n", cfg)

	if err := http.ListenAndServe(cfg.address, mux); err != nil {
		fmt.Printf("Listening ends with error %s", err.Error())
	}
}
