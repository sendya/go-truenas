package truenas

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	t.Parallel()
	t.Run("Error message", func(t *testing.T) {
		err := NewNotFoundError("pool", "ID 123")
		assert.Equal(t, "pool with ID 123 not found", err.Error())
	})

	t.Run("Is method", func(t *testing.T) {
		err1 := NewNotFoundError("pool", "ID 123")
		err2 := NewNotFoundError("dataset", "name test")
		target := &NotFoundError{}

		assert.True(t, errors.Is(err1, target))
		assert.True(t, errors.Is(err2, target))

		var extracted *NotFoundError
		assert.True(t, errors.As(err1, &extracted))
		assert.Equal(t, "pool", extracted.ResourceType)
	})

	t.Run("Different error type", func(t *testing.T) {
		err := NewNotFoundError("pool", "ID 123")
		otherErr := errors.New("some other error")

		assert.False(t, errors.Is(otherErr, err))
	})
}

func TestNotFoundErrorUsage(t *testing.T) {
	t.Parallel()
	var err error = NewNotFoundError("pool", "ID 123")

	// Test direct type assertion
	if notFoundErr, ok := err.(*NotFoundError); ok {
		assert.Equal(t, "pool", notFoundErr.ResourceType)
		assert.Equal(t, "ID 123", notFoundErr.Identifier)
	}

	// Test errors.Is usage
	target := &NotFoundError{}
	assert.True(t, errors.Is(err, target))

	// Test errors.As usage
	var extracted *NotFoundError
	if errors.As(err, &extracted) {
		assert.Equal(t, "pool", extracted.ResourceType)
		assert.Equal(t, "ID 123", extracted.Identifier)
	}
}
