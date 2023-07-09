package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/metric/server"
	"github.com/oriolus/notprometheus/internal/metric/storage"
)

var (
	badUriFormatError    = errors.New("bad uri format")
	notImplementedMetric = errors.New("unknown type of metric")
)

const methodName = "update"

type UpdateHandler struct {
	server *server.Server
}

func NewUpdateHandler(server *server.Server) *UpdateHandler {
	return &UpdateHandler{server: server}
}

func (s *UpdateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	method, typ, name, value, err := parseRequestURI(req.RequestURI)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if method != methodName {
		res.WriteHeader(http.StatusNotFound)
		return
	}

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
		return s.processCounter(name)
	} else if mType == metric.TypeGauge {
		return s.processGauge(name, value)
	} else {
		return notImplementedMetric
	}
}

func (s *UpdateHandler) processCounter(name string) error {
	c, err := s.server.Storage().GetCounter(name)
	if err == nil {
		c.Inc()
		return nil
	}

	if !errors.Is(err, storage.MetricNotFound) {
		return err
	}

	c, err = metric.NewCounter(name)
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

	if !errors.Is(err, storage.MetricNotFound) {
		return err
	}

	g, err = metric.NewGauge(name, v)
	if err != nil {
		return err
	}
	return s.server.Storage().SetGauge(g)
}

func parseRequestURI(requestUri string) (method, typ, name, value string, err error) {
	parts := strings.Split(requestUri, "/")
	if len(parts) != 5 {
		return "", "", "", "", badUriFormatError
	}

	method = parts[1]
	typ = parts[2]
	name = parts[3]
	value = parts[4]
	err = nil

	return
}
