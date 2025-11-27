package assert_test

import (
	"testing"

	"github.com/jespert/primordial/internal/quality/assert"
)

func TestEqual(t *testing.T) {
	t.Run("Do not panic if equal", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.Equal(true, true)
		})
	})
	t.Run("Panic if different", func(t *testing.T) {
		expectPanic(t, "expected false == true", func() {
			assert.Equal(false, true)
		})
	})
}

func TestEqualf(t *testing.T) {
	t.Run("Do not panic if equal", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.Equalf(true, true, format, value)
		})
	})
	t.Run("Panic if different", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assert.Equalf(false, true, format, value)
		})
	})
}
