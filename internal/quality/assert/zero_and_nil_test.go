package assert_test

import (
	"testing"

	"github.com/jespert/primordial/internal/quality/assert"
)

func TestNotZero(t *testing.T) {
	t.Run("Do not panic if not zero", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.NotZero(new(int))
		})
	})
	t.Run("Panic if zero", func(t *testing.T) {
		expectPanic(t, "zero value", func() {
			assert.NotZero((*int)(nil))
		})
	})
}

func TestNotZerof(t *testing.T) {
	t.Run("Do not panic if not zero", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.NotZerof(new(int), format, value)
		})
	})
	t.Run("Panic if zero", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assert.NotZerof((*int)(nil), format, value)
		})
	})
}

func TestNotNilSlice(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.NotNilSlice([]int{})
		})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		expectPanic(t, "nil slice", func() {
			assert.NotNilSlice([]int(nil))
		})
	})
}

func TestNotNilSlicef(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.NotNilSlicef([]int{}, format, value)
		})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assert.NotNilSlicef([]int(nil), format, value)
		})
	})
}

func TestNotNilMap(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.NotNilMap(map[bool]int{})
		})
	})
	t.Run("Panic if zero", func(t *testing.T) {
		expectPanic(t, "nil map", func() {
			assert.NotNilMap(map[bool]int(nil))
		})
	})
}

func TestNotNilMapf(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assert.NotNilMapf(map[bool]int{}, format, value)
		})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assert.NotNilMapf(map[bool]int(nil), format, value)
		})
	})
}
