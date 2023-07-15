package http

import (
	"fmt"
	"net/http"

	"github.com/oriolus/notprometheus/internal/collector/sender"
	"github.com/oriolus/notprometheus/internal/metric"
)

var _ sender.MetricSender = (*Client)(nil)

type Client struct {
	client http.Client
	base   string
}

func NewClient(base string) (*Client, error) {
	if base == "" {
		return nil, sender.ErrStringIsEmpty
	}
	return &Client{client: http.Client{}, base: base}, nil
}

func (c *Client) UpdateGauge(gauge metric.Gauge) error {
	url := c.getURL(metric.TypeGauge, gauge.Name(), fmt.Sprintf("%f", gauge.Value()))
	resp, err := c.client.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) UpdateCounter(counter metric.Counter) error {
	url := c.getURL(metric.TypeCounter, counter.Name(), fmt.Sprintf("%d", counter.Value()))
	resp, err := c.client.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) getURL(metricType metric.Type, name, value string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s", c.base, string(metricType), name, value)
}
