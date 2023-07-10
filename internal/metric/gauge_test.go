package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGauge(t *testing.T) {
	t.Parallel()

	t.Run("NewGauge; name is empty; should return error", func(t *testing.T) {
		t.Parallel()
		// arrange
		name := ""

		// act
		_, err := NewGauge(name, 1.0)

		// assert
		require.Error(t, err)
		assert.Equal(t, EmptyNameError, err)
	})

	t.Run("NewGauge; name is not empty; should not return error", func(t *testing.T) {
		t.Parallel()
		// arrange
		name := "name"

		// act
		g, err := NewGauge(name, 1.0)

		// assert
		require.NoError(t, err)
		assert.NotNil(t, g)
	})
}
