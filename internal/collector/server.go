package collector

import (
	"context"
	"time"
)

type Server struct {
}

func (s *Server) Run(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := processMetrics(ctx)
			if err != nil {
				// general have to process error
				// but for now its enough
				return err
			}
		}
	}
}

func processMetrics(ctx context.Context) error {
	panic("not implemented")
}
