package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/server"
)

type GetMetricHandler struct {
	server *server.Server
}

func NewGetMetricHandler(server *server.Server) (*GetMetricHandler, error) {
	if server == nil {
		return nil, ArgumentNilError
	}

	return &GetMetricHandler{server: server}, nil
}

func (s *GetMetricHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, URLParamMetricType)
	metricName := chi.URLParam(req, URLParamName)

	mType, err := metric.GetMetricType(metricType)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if mType == metric.TypeGauge {
		gauge, err := s.server.Storage().GetGauge(metricName)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = res.Write([]byte(fmt.Sprintf("%f", gauge.Value())))
		if err != nil {
			fmt.Printf("errorw while writing to response %s\r\n", err.Error())
		}
	} else if mType == metric.TypeCounter {
		cnt, err := s.server.Storage().GetCounter(metricName)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = res.Write([]byte(fmt.Sprintf("%d", cnt.Value())))
		if err != nil {
			fmt.Printf("errorw while writing to response %s\r\n", err.Error())
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
