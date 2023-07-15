package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage"
)

var (
	ErrBadURIFormat   = errors.New("bad uri format")
	ErrNotImplemented = errors.New("unknown type of metric")
)

type UpdateHandler struct {
	server *server.Server
}

func NewUpdateHandler(server *server.Server) (*UpdateHandler, error) {
	if server == nil {
		return nil, ArgumentNilError
	}
	return &UpdateHandler{server: server}, nil
}

func (s *UpdateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	typ := chi.URLParam(req, URLParamMetricType)
	name := chi.URLParam(req, URLParamName)
	value := chi.URLParam(req, URLParamValue)

	mType, err := metric.GetMetricType(typ)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.handle(mType, name, value)
	if err != nil {
		// todo: differentiate user errors and server errors
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (s *UpdateHandler) handle(mType metric.Type, name string, value string) error {
	if mType == metric.TypeCounter {
		return s.processCounter(name, value)
	} else if mType == metric.TypeGauge {
		return s.processGauge(name, value)
	} else {
		return ErrNotImplemented
	}
}

func (s *UpdateHandler) processCounter(name, value string) error {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	c, err := s.server.Storage().GetCounter(name)
	if err == nil {
		c.Set(intValue)
		return nil
	}

	if err != storage.ErrMetricNotFound {
		return err
	}

	c, err = metric.NewCounter(name)
	if err != nil {
		return err
	}
	c, err = metric.NewCounterWithValue(name, intValue)
	if err != nil {
		return err
	}
	return s.server.Storage().SetCounter(c)
}

func (s *UpdateHandler) processGauge(name, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("error while parsing gauge value %s: %s", value, err)
	}

	g, err := s.server.Storage().GetGauge(name)
	if err == nil {
		g.Set(v)
		return nil
	}

	if !errors.Is(err, storage.ErrMetricNotFound) {
		return err
	}

	g, err = metric.NewGauge(name, v)
	if err != nil {
		return err
	}
	return s.server.Storage().SetGauge(g)
}
