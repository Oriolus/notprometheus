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

type UpdateJSONedHandler struct {
	server *server.Server
}

func NewUpdateJSONedHandler(server *server.Server) (*UpdateJSONedHandler, error) {
	if server == nil {
		return nil, ErrArgumentNil
	}
	return &UpdateJSONedHandler{server: server}, nil
}

func (s *UpdateJSONedHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add(contentTypeKey, JSONContentTypeValue)

	if req.Header.Get(contentTypeKey) != JSONContentTypeValue {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(req.Body); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var metricRequest models.UpdateMetricRequest
	if err := json.Unmarshal(buffer.Bytes(), &metricRequest); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validate(metricRequest); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.handle(metricRequest)
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

	res.WriteHeader(http.StatusOK)
}

func validate(req models.UpdateMetricRequest) error {
	mType := metric.Type(req.MType)
	switch mType {
	case metric.TypeCounter:
		if req.Delta == nil {
			return errInvalidCounterReq
		}
	case metric.TypeGauge:
		if req.Value == nil {
			return errInvalidGaugeReq
		}
	default:
		return ErrNotImplemented
	}

	return nil
}

func (s *UpdateJSONedHandler) handle(req models.UpdateMetricRequest) (*models.UpdateMetricResponse, error) {
	metricType := metric.Type(req.MType)
	switch metricType {
	case metric.TypeCounter:
		{
			c, err := s.processCounter(req.ID, *req.Delta)
			if err != nil {
				return nil, err
			}

			return &models.UpdateMetricResponse{
				ID:    c.Name(),
				MType: req.MType,
				Value: float64(c.Value()),
			}, nil
		}
	case metric.TypeGauge:
		{
			g, err := s.processGauge(req.ID, *req.Value)
			if err != nil {
				return nil, err
			}

			return &models.UpdateMetricResponse{
				ID:    g.Name(),
				MType: req.MType,
				Value: g.Value(),
			}, nil
		}
	default:
		panic("unexpected metric type")
	}
}

func (s *UpdateJSONedHandler) processCounter(name string, value int64) (metric.Counter, error) {
	c, err := s.server.Storage().GetCounter(name)
	if err == nil {
		c.Add(value)
		return c, nil
	}

	if err != storage.ErrMetricNotFound {
		return nil, err
	}

	c, err = metric.NewCounterWithValue(name, value)
	if err != nil {
		return nil, err
	}
	return c, s.server.Storage().SetCounter(c)
}

func (s *UpdateJSONedHandler) processGauge(name string, value float64) (metric.Gauge, error) {
	g, err := s.server.Storage().GetGauge(name)
	if err == nil {
		g.Set(value)
		return g, nil
	}

	if !errors.Is(err, storage.ErrMetricNotFound) {
		return nil, err
	}

	g, err = metric.NewGauge(name, value)
	if err != nil {
		return nil, err
	}
	return g, s.server.Storage().SetGauge(g)
}
