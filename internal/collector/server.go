package collector

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/oriolus/notprometheus/internal/collector/sender"
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

	pollCount, err := metric.NewCounter("poolCount")
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

	alloc, err := metric.NewGauge("memstat:Alloc", float64(memStat.Alloc))
	if err == nil {
		gauges = append(gauges, alloc)
	}
	backHashSys, err := metric.NewGauge("memstat:BuckHashSys", float64(memStat.BuckHashSys))
	if err == nil {
		gauges = append(gauges, backHashSys)
	}
	frees, err := metric.NewGauge("memstat:Frees", float64(memStat.Frees))
	if err == nil {
		gauges = append(gauges, frees)
	}
	gCCpuFraction, err := metric.NewGauge("memestat:GCCPUFraction", memStat.GCCPUFraction)
	if err == nil {
		gauges = append(gauges, gCCpuFraction)
	}
	gCSys, err := metric.NewGauge("memStat:GCSys", float64(memStat.GCSys))
	if err == nil {
		gauges = append(gauges, gCSys)
	}
	heapAlloc, err := metric.NewGauge("memStat:HeapAlloc", float64(memStat.HeapAlloc))
	if err == nil {
		gauges = append(gauges, heapAlloc)
	}
	heapIdle, err := metric.NewGauge("memStat:HeapIdle", float64(memStat.HeapIdle))
	if err == nil {
		gauges = append(gauges, heapIdle)
	}
	heapInuse, err := metric.NewGauge("memStat:HeapInuse", float64(memStat.HeapInuse))
	if err == nil {
		gauges = append(gauges, heapInuse)
	}
	heapObjects, err := metric.NewGauge("memStat:HeapObjects", float64(memStat.HeapObjects))
	if err == nil {
		gauges = append(gauges, heapObjects)
	}
	heapReleased, err := metric.NewGauge("memStat:HeapReleased", float64(memStat.HeapReleased))
	if err == nil {
		gauges = append(gauges, heapReleased)
	}
	heapSys, err := metric.NewGauge("memStat:HeapSys", float64(memStat.HeapSys))
	if err == nil {
		gauges = append(gauges, heapSys)
	}
	lastGC, err := metric.NewGauge("memStat:LastGC", float64(memStat.LastGC))
	if err == nil {
		gauges = append(gauges, lastGC)
	}
	lookups, err := metric.NewGauge("memStat:Lookups", float64(memStat.Lookups))
	if err == nil {
		gauges = append(gauges, lookups)
	}
	mCacheInuse, err := metric.NewGauge("memStat:MCacheInuse", float64(memStat.MCacheInuse))
	if err == nil {
		gauges = append(gauges, mCacheInuse)
	}
	mCacheSys, err := metric.NewGauge("memStat:MCacheSys", float64(memStat.MCacheSys))
	if err == nil {
		gauges = append(gauges, mCacheSys)
	}
	mSpanInuse, err := metric.NewGauge("memStat:MSpanInuse", float64(memStat.MSpanInuse))
	if err == nil {
		gauges = append(gauges, mSpanInuse)
	}
	mSpanSys, err := metric.NewGauge("memStat:MSpanSys", float64(memStat.MSpanSys))
	if err == nil {
		gauges = append(gauges, mSpanSys)
	}
	mallocs, err := metric.NewGauge("memStat:Mallocs", float64(memStat.Mallocs))
	if err == nil {
		gauges = append(gauges, mallocs)
	}
	nextGC, err := metric.NewGauge("memStat:NextGC", float64(memStat.NextGC))
	if err == nil {
		gauges = append(gauges, nextGC)
	}
	numForcedGC, err := metric.NewGauge("memStat:NumForcedGC", float64(memStat.NumForcedGC))
	if err == nil {
		gauges = append(gauges, numForcedGC)
	}
	numGC, err := metric.NewGauge("memStat:NumGC", float64(memStat.NumGC))
	if err == nil {
		gauges = append(gauges, numGC)
	}
	otherSys, err := metric.NewGauge("memStat:OtherSys", float64(memStat.OtherSys))
	if err == nil {
		gauges = append(gauges, otherSys)
	}
	pauseTotalNs, err := metric.NewGauge("memStat:PauseTotalNs", float64(memStat.PauseTotalNs))
	if err == nil {
		gauges = append(gauges, pauseTotalNs)
	}
	stackInuse, err := metric.NewGauge("memStat:StackInuse", float64(memStat.StackInuse))
	if err == nil {
		gauges = append(gauges, stackInuse)
	}
	stackSys, err := metric.NewGauge("memStat:StackSys", float64(memStat.StackSys))
	if err == nil {
		gauges = append(gauges, stackSys)
	}
	sys, err := metric.NewGauge("memStat:Sys", float64(memStat.Sys))
	if err == nil {
		gauges = append(gauges, sys)
	}
	totalAlloc, err := metric.NewGauge("memStat:TotalAlloc", float64(memStat.TotalAlloc))
	if err == nil {
		gauges = append(gauges, totalAlloc)
	}
	random, err := metric.NewGauge("random", rand.Float64())
	if err == nil {
		gauges = append(gauges, random)
	}

	s.gauges = gauges
	s.pollCount.Inc()

	return nil
}

func (s *Server) send() {
	var err error
	for _, g := range s.gauges {
		err = s.client.UpdateGauge(g)
		if err != nil {
			fmt.Printf("metric %s was not updated due %s", g.Name(), err.Error())
		}
	}

	err = s.client.UpdateCounter(s.pollCount)
	if err != nil {
		fmt.Printf("metric %s was not updated due %s", s.pollCount.Name(), err.Error())
	}
}
