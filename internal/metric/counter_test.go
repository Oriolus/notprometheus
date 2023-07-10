package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	t.Parallel()

	t.Run("NewCounter; name is empty; should return error", func(t *testing.T) {
		t.Parallel()
		// arrange
		name := ""

		// act
		_, err := NewCounter(name)

		// assert
		require.Error(t, err)
		assert.Equal(t, EmptyNameError, err)
	})

	t.Run("NewCounter; name is not empty; should not return error", func(t *testing.T) {
		t.Parallel()
		// arrange
		name := "name"

		// act
		cnt, err := NewCounter(name)

		// assert
		require.NoError(t, err)
		assert.NotNil(t, cnt)
	})
}
