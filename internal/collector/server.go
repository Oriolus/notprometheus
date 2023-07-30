package collector

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/oriolus/notprometheus/internal/collector/sender"
	"github.com/oriolus/notprometheus/internal/logger"
	"github.com/oriolus/notprometheus/internal/metric"
)

type Server struct {
	client         sender.MetricSender
	pollCount      metric.Counter
	gauges         []metric.Gauge
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewServer(client sender.MetricSender, pollInterval, reportInterval time.Duration) (*Server, error) {
	if client == nil {
		return nil, ErrArgumentNil
	}

	pollCount, err := metric.NewCounter("PoolCount")
	if err != nil {
		return nil, err
	}
	return &Server{
		client:         client,
		pollCount:      pollCount,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	pollTicker := time.NewTicker(s.pollInterval)
	reportTicker := time.NewTicker(s.reportInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-pollTicker.C:
			err := s.processMetrics()
			if err != nil {
				fmt.Println(err)
			}
		case <-reportTicker.C:
			s.send()
		}
	}
}

func (s *Server) processMetrics() error {
	memStat := runtime.MemStats{}
	runtime.ReadMemStats(&memStat)

	gauges := make([]metric.Gauge, 0)

	alloc, err := metric.NewGauge("Alloc", float64(memStat.Alloc))
	if err == nil {
		gauges = append(gauges, alloc)
	}
	backHashSys, err := metric.NewGauge("BuckHashSys", float64(memStat.BuckHashSys))
	if err == nil {
		gauges = append(gauges, backHashSys)
	}
	frees, err := metric.NewGauge("Frees", float64(memStat.Frees))
	if err == nil {
		gauges = append(gauges, frees)
	}
	gCCpuFraction, err := metric.NewGauge("GCCPUFraction", memStat.GCCPUFraction)
	if err == nil {
		gauges = append(gauges, gCCpuFraction)
	}
	gCSys, err := metric.NewGauge("GCSys", float64(memStat.GCSys))
	if err == nil {
		gauges = append(gauges, gCSys)
	}
	heapAlloc, err := metric.NewGauge("HeapAlloc", float64(memStat.HeapAlloc))
	if err == nil {
		gauges = append(gauges, heapAlloc)
	}
	heapIdle, err := metric.NewGauge("HeapIdle", float64(memStat.HeapIdle))
	if err == nil {
		gauges = append(gauges, heapIdle)
	}
	heapInuse, err := metric.NewGauge("HeapInuse", float64(memStat.HeapInuse))
	if err == nil {
		gauges = append(gauges, heapInuse)
	}
	heapObjects, err := metric.NewGauge("HeapObjects", float64(memStat.HeapObjects))
	if err == nil {
		gauges = append(gauges, heapObjects)
	}
	heapReleased, err := metric.NewGauge("HeapReleased", float64(memStat.HeapReleased))
	if err == nil {
		gauges = append(gauges, heapReleased)
	}
	heapSys, err := metric.NewGauge("HeapSys", float64(memStat.HeapSys))
	if err == nil {
		gauges = append(gauges, heapSys)
	}
	lastGC, err := metric.NewGauge("LastGC", float64(memStat.LastGC))
	if err == nil {
		gauges = append(gauges, lastGC)
	}
	lookups, err := metric.NewGauge("Lookups", float64(memStat.Lookups))
	if err == nil {
		gauges = append(gauges, lookups)
	}
	mCacheInuse, err := metric.NewGauge("MCacheInuse", float64(memStat.MCacheInuse))
	if err == nil {
		gauges = append(gauges, mCacheInuse)
	}
	mCacheSys, err := metric.NewGauge("MCacheSys", float64(memStat.MCacheSys))
	if err == nil {
		gauges = append(gauges, mCacheSys)
	}
	mSpanInuse, err := metric.NewGauge("MSpanInuse", float64(memStat.MSpanInuse))
	if err == nil {
		gauges = append(gauges, mSpanInuse)
	}
	mSpanSys, err := metric.NewGauge("MSpanSys", float64(memStat.MSpanSys))
	if err == nil {
		gauges = append(gauges, mSpanSys)
	}
	mallocs, err := metric.NewGauge("Mallocs", float64(memStat.Mallocs))
	if err == nil {
		gauges = append(gauges, mallocs)
	}
	nextGC, err := metric.NewGauge("NextGC", float64(memStat.NextGC))
	if err == nil {
		gauges = append(gauges, nextGC)
	}
	numForcedGC, err := metric.NewGauge("NumForcedGC", float64(memStat.NumForcedGC))
	if err == nil {
		gauges = append(gauges, numForcedGC)
	}
	numGC, err := metric.NewGauge("NumGC", float64(memStat.NumGC))
	if err == nil {
		gauges = append(gauges, numGC)
	}
	otherSys, err := metric.NewGauge("OtherSys", float64(memStat.OtherSys))
	if err == nil {
		gauges = append(gauges, otherSys)
	}
	pauseTotalNs, err := metric.NewGauge("PauseTotalNs", float64(memStat.PauseTotalNs))
	if err == nil {
		gauges = append(gauges, pauseTotalNs)
	}
	stackInuse, err := metric.NewGauge("StackInuse", float64(memStat.StackInuse))
	if err == nil {
		gauges = append(gauges, stackInuse)
	}
	stackSys, err := metric.NewGauge("StackSys", float64(memStat.StackSys))
	if err == nil {
		gauges = append(gauges, stackSys)
	}
	sys, err := metric.NewGauge("Sys", float64(memStat.Sys))
	if err == nil {
		gauges = append(gauges, sys)
	}
	totalAlloc, err := metric.NewGauge("TotalAlloc", float64(memStat.TotalAlloc))
	if err == nil {
		gauges = append(gauges, totalAlloc)
	}
	random, err := metric.NewGauge("RandomValue", rand.Float64())
	if err == nil {
		gauges = append(gauges, random)
	}

	s.gauges = gauges
	s.pollCount.Inc()

	return nil
}

func (s *Server) send() {
	var err error
	//for _, g := range s.gauges {
	//	err = s.client.UpdateGauge(g)
	//	if err != nil {
	//		logger.Errorf("metric %s was not updated due %s", g.Name(), err.Error())
	//	}
	//}

	err = s.client.UpdateCounter(s.pollCount)
	if err != nil {
		logger.Errorf("metric %s was not updated due %s", s.pollCount.Name(), err.Error())
	}
}
