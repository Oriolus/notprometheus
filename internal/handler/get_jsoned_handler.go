package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/models"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage"
)

type GetJSONedHandler struct {
	server *server.Server
}

func NewGetJSONedHandler(server *server.Server) (*GetJSONedHandler, error) {
	if server == nil {
		return nil, ErrArgumentNil
	}
	return &GetJSONedHandler{server: server}, nil
}

func (s *GetJSONedHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get(contentTypeKey) != JSONContentTypeValue {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(req.Body); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var metricRequest models.GetMetricRequest
	if err := json.Unmarshal(buffer.Bytes(), &metricRequest); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isMetricTypeValid(metricRequest.MType) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.handle(metricRequest)
	if errors.Is(err, storage.ErrMetricNotFound) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBuff, err := json.Marshal(response)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = res.Write(respBuff); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(contentTypeKey, JSONContentTypeValue)
	res.WriteHeader(http.StatusOK)
}

func isMetricTypeValid(metricType string) bool {
	_, err := metric.GetMetricType(metricType)
	return err == nil
}

func (s *GetJSONedHandler) handle(req models.GetMetricRequest) (*models.GetMetricResponse, error) {
	metricType := metric.Type(req.MType)
	switch metricType {
	case metric.TypeCounter:
		{
			c, err := s.server.Storage().GetCounter(req.ID)
			if err != nil {
				return nil, err
			}

			return &models.GetMetricResponse{
				ID:    c.Name(),
				MType: req.MType,
				Value: float64(c.Value()),
			}, nil
		}
	case metric.TypeGauge:
		{
			g, err := s.server.Storage().GetGauge(req.ID)
			if err != nil {
				return nil, err
			}

			return &models.GetMetricResponse{
				ID:    g.Name(),
				MType: req.MType,
				Value: g.Value(),
			}, nil
		}
	default:
		panic("unexpected metric type")
	}
}
