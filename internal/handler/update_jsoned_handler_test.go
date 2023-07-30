package handler

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/models"
	"github.com/oriolus/notprometheus/internal/server"
	"github.com/oriolus/notprometheus/internal/server/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newIDGenerator() func() int64 {
	seed := atomic.Int64{}
	seed.Store(time.Now().Unix())
	return func() int64 {
		seed.Add(1)
		return seed.Load()
	}
}

var newID = newIDGenerator()

func newMetricID() string {
	ID := newID()
	return fmt.Sprintf("order number: %d", ID)
}

func TestValidate(t *testing.T) {
	t.Parallel()

	getIntPointer := func() *int64 {
		val := int64(2)
		return &val
	}

	getFloatPointer := func() *float64 {
		val := 2.2
		return &val
	}

	testCases := []struct {
		name     string
		req      models.UpdateMetricRequest
		expected error
	}{
		{"validate; empty type; should return not implemented error", models.UpdateMetricRequest{MType: ""}, ErrNotImplemented},
		{"validate; gauge with empty Value; should return error", models.UpdateMetricRequest{MType: string(metric.TypeGauge)}, errInvalidGaugeReq},
		{"validate; counter with empty Delta; should return error", models.UpdateMetricRequest{MType: string(metric.TypeCounter)}, errInvalidCounterReq},
		{"validate; valid gauge; should return no error", models.UpdateMetricRequest{MType: string(metric.TypeGauge), Value: getFloatPointer()}, nil},
		{"validate; valid counter; should return no error", models.UpdateMetricRequest{MType: string(metric.TypeCounter), Delta: getIntPointer()}, nil},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			// arrange
			t.Parallel()

			// act
			err := validate(testCase.req)

			// assert
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestHandle(t *testing.T) {
	t.Parallel()

	storage := memory.NewMemStorage()
	memServer := server.NewServer(storage)
	handler, _ := NewUpdateJSONedHandler(memServer)

	t.Run("handle; new counter; should return new counter", func(t *testing.T) {
		// arrange
		t.Parallel()

		delta := int64(10)
		req := models.UpdateMetricRequest{
			ID:    newMetricID(),
			MType: string(metric.TypeCounter),
			Delta: &delta,
		}
		expectedResp := &models.UpdateMetricResponse{
			ID:    req.ID,
			MType: req.MType,
			Delta: &delta,
		}

		// act
		resp, err := handler.handle(req)

		// assert
		require.NoError(t, err)
		assert.Equal(t, expectedResp, resp)

		cnt, err := storage.GetCounter(req.ID)
		require.NoError(t, err)
		assert.Equal(t, delta, cnt.Value())
	})

	t.Run("handle; exists counter; should return updated counter", func(t *testing.T) {
		// arrange
		t.Parallel()

		seed := int64(10)
		cnt, err := metric.NewCounterWithValue(newMetricID(), seed)
		require.NoError(t, err)

		err = storage.SetCounter(cnt)
		require.NoError(t, err)

		delta := int64(3)
		req := models.UpdateMetricRequest{
			ID:    cnt.Name(),
			MType: string(metric.TypeCounter),
			Delta: &delta,
		}
		res := seed + delta
		expectedResp := &models.UpdateMetricResponse{
			ID:    cnt.Name(),
			MType: req.MType,
			Delta: &res,
		}
		expectedErr := error(nil)

		// act
		resp, err := handler.handle(req)

		// assert
		require.Equal(t, expectedErr, err)
		assert.Equal(t, expectedResp, resp)

		actualCnt, err := storage.GetCounter(req.ID)
		require.NoError(t, err)
		assert.Equal(t, seed+delta, actualCnt.Value())
	})

	t.Run("handle; new gauge; should return gauge", func(t *testing.T) {
		// arrange
		t.Parallel()

		val := 2.2
		req := models.UpdateMetricRequest{
			ID:    newMetricID(),
			MType: string(metric.TypeGauge),
			Value: &val,
		}
		expectedResp := &models.UpdateMetricResponse{
			ID:    req.ID,
			MType: req.MType,
			Value: req.Value,
		}
		expectedErr := error(nil)

		// act
		resp, err := handler.handle(req)

		// assert
		require.Equal(t, expectedErr, err)
		assert.Equal(t, expectedResp, resp)

		actualGauge, err := storage.GetGauge(req.ID)
		require.NoError(t, err)
		assert.Equal(t, actualGauge.Value(), val)
	})

	t.Run("handle; exists gauge; should return updated gauge", func(t *testing.T) {
		// arrange
		t.Parallel()

		seed := 3.3
		g, err := metric.NewGauge(newMetricID(), seed)
		require.NoError(t, err)

		err = storage.SetGauge(g)
		require.NoError(t, err)

		newValue := 5.5
		req := models.UpdateMetricRequest{
			ID:    g.Name(),
			MType: string(metric.TypeGauge),
			Value: &newValue,
		}
		expectedResp := &models.UpdateMetricResponse{
			ID:    req.ID,
			MType: req.MType,
			Value: &newValue,
		}
		expectedErr := error(nil)

		// act
		resp, err := handler.handle(req)

		// assert
		require.Equal(t, expectedErr, err)
		assert.Equal(t, expectedResp, resp)

		actualGauge, err := storage.GetGauge(g.Name())
		require.NoError(t, err)
		assert.Equal(t, actualGauge.Value(), newValue)
	})
}
