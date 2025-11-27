package assert_test

import (
	"errors"
	"testing"

	"github.com/jespert/primordial/internal/quality/assert"
)

func TestSuccess(t *testing.T) {
	t.Run("Do not panic if nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.Success(nil)
		})
	})
	t.Run("Panic if not nil", func(t *testing.T) {
		expectPanic(t, "unexpected error: test error", func() {
			assert.Success(errors.New("test error"))
		})
	})
}

func TestSuccessf(t *testing.T) {
	t.Run("Do not panic if nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.Successf(nil, format, value)
		})
	})
	t.Run("Panic if not nil", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			err := errors.New("test error")
			assert.Successf(err, format, value)
		})
	})
}
