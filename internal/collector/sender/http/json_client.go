package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oriolus/notprometheus/internal/collector/sender"
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/models"
)

var _ sender.MetricSender = (*jsonClient)(nil)

const (
	JSONContentTypeValue = "application/json"
)

type jsonClient struct {
	client http.Client
	base   string
}

func NewJSONClient(base string) (sender.MetricSender, error) {
	if base == "" {
		return nil, sender.ErrStringIsEmpty
	}
	return &jsonClient{client: http.Client{Timeout: sender.ClientTimeout}, base: base}, nil
}

func (c *jsonClient) UpdateGauge(gauge metric.Gauge) error {
	url := c.getURL()
	val := gauge.Value()
	req := models.UpdateMetricRequest{
		ID:    gauge.Name(),
		MType: string(metric.TypeGauge),
		Value: &val,
	}

	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(url, JSONContentTypeValue, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *jsonClient) UpdateCounter(counter metric.Counter) error {
	url := c.getURL()
	delta := counter.Value()
	req := models.UpdateMetricRequest{
		ID:    counter.Name(),
		MType: string(metric.TypeCounter),
		Delta: &delta,
	}

	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(url, JSONContentTypeValue, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *jsonClient) getURL() string {
	return fmt.Sprintf("%s/update", c.base)
}
