package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/metric/server"
)

var (
	badUriFormatError = errors.New("bad uri format")
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

	s.handle(res, mType, name, value)
}

func (s *UpdateHandler) handle(res http.ResponseWriter, p metric.Type, name string, value string) error {
	res.WriteHeader(http.StatusOK)
	return nil
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
