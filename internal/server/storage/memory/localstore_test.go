package memory

import (
	"testing"

	"github.com/oriolus/notprometheus/internal/metric"
	storageImp "github.com/oriolus/notprometheus/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage(t *testing.T) {
	t.Parallel()

	t.Run("GetGauge; metric not found; should return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		name := "name"
		storage := NewMemStorage()

		// act
		_, err := storage.GetGauge(name)

		// assert
		require.Error(t, err)
		assert.Equal(t, err, storageImp.MetricNotFoundError)
	})

	t.Run("GetGauge; metric found; should not return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		g, err := metric.NewGauge("gauge", 1.0)
		require.NoError(t, err)
		storage := NewMemStorage()
		err = storage.SetGauge(g)
		require.NoError(t, err)

		// act
		actualG, err := storage.GetGauge(g.Name())

		// assert
		require.NoError(t, err)
		assert.NotNil(t, actualG)
	})

	t.Run("SetGauge; metric is nil; should return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		storage := NewMemStorage()

		// act
		err := storage.SetGauge(nil)

		// assert
		require.Error(t, err)
		assert.Error(t, err, storageImp.ArgumentNilError)
	})

	t.Run("SetGauge; metric is nil; should return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		storage := NewMemStorage()
		g, err := metric.NewGauge("name", 1.0)
		require.NoError(t, err)

		// act
		err = storage.SetGauge(g)

		// assert
		require.NoError(t, err)
	})

	t.Run("GetCounter; metric not found; should return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		name := "name"
		storage := NewMemStorage()

		// act
		_, err := storage.GetCounter(name)

		// assert
		require.Error(t, err)
		assert.Equal(t, err, storageImp.MetricNotFoundError)
	})

	t.Run("GetCounter; metric found; should not return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		c, err := metric.NewCounter("counter")
		require.NoError(t, err)
		storage := NewMemStorage()
		err = storage.SetCounter(c)
		require.NoError(t, err)

		// act
		actualC, err := storage.GetCounter(c.Name())

		// assert
		require.NoError(t, err)
		assert.NotNil(t, actualC)
	})

	t.Run("SetCounter; metric is nil; should return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		storage := NewMemStorage()

		// act
		err := storage.SetCounter(nil)

		// assert
		require.Error(t, err)
		assert.Error(t, err, storageImp.ArgumentNilError)
	})

	t.Run("SetCounter; metric is nil; should return error;", func(t *testing.T) {
		t.Parallel()
		// arrange
		storage := NewMemStorage()
		c, err := metric.NewCounter("name")
		require.NoError(t, err)

		// act
		err = storage.SetCounter(c)

		// assert
		require.NoError(t, err)
	})
}
