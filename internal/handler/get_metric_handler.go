package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/server"
	"net/http"
	"strconv"
)

type GetMetricHandler struct {
	server *server.Server
}

func NewGetMetricHandler(server *server.Server) (*GetMetricHandler, error) {
	if server == nil {
		return nil, ErrArgumentNil
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

	switch mType {
	case metric.TypeGauge:
		{
			gauge, err := s.server.Storage().GetGauge(metricName)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}

			_, err = res.Write([]byte(strconv.FormatFloat(gauge.Value(), 'f', -1, 64)))
			if err != nil {
				// вообще это 500, но все же)
				fmt.Printf("errorw while writing to response %s\r\n", err.Error())
			}
		}
	case metric.TypeCounter:
		{
			cnt, err := s.server.Storage().GetCounter(metricName)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}

			_, err = res.Write([]byte(fmt.Sprintf("%d", cnt.Value())))
			if err != nil {
				fmt.Printf("errorw while writing to response %s\r\n", err.Error())
			}
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}
