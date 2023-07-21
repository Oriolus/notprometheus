package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/server"
)

type GetAllHandler struct {
	server *server.Server
}

func NewGetAllHandler(server *server.Server) (*GetAllHandler, error) {
	if server == nil {
		return nil, ErrArgumentNil
	}

	return &GetAllHandler{server: server}, nil
}

func (s *GetAllHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	builder := strings.Builder{}

	for _, g := range s.server.Storage().AllGauges() {
		builder.WriteString(getGauge(g))
		builder.WriteString("\n")
	}

	for _, c := range s.server.Storage().AllCounters() {
		builder.WriteString(getCounter(c))
		builder.WriteString("\n")
	}

	_, err := res.Write([]byte(builder.String()))
	if err != nil {
		fmt.Println(err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func getGauge(g metric.Gauge) string {
	return fmt.Sprintf("%s: %f", g.Name(), g.Value())
}

func getCounter(c metric.Counter) string {
	return fmt.Sprintf("%s: %d", c.Name(), c.Value())
}
