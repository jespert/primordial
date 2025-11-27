package assert_test

import (
	"testing"

	"github.com/jespert/primordial/internal/quality/assert"
)

func TestTrue(t *testing.T) {
	t.Run("Do not panic if true", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.True(true)
		})
	})
	t.Run("Panic if false", func(t *testing.T) {
		expectPanic(t, "expected condition to be true", func() {
			assert.True(false)
		})
	})
}

func TestTruef(t *testing.T) {
	t.Run("Do not panic if true", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.Truef(true, format, value)
		})
	})
	t.Run("Panic if false", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assert.Truef(false, format, value)
		})
	})
}
