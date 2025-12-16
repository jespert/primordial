// Package require provides assertions for testing that abort the test on failure.
package require

import (
	"github.com/jespert/primordial/internal/quality/expect"
)

type TestingT interface {
	expect.TestingT
	FailNow()
}

func Equal[T comparable](t TestingT, want, got T) {
	t.Helper()
	if !expect.Equal(t, want, got) {
		t.FailNow()
	}
}

func Success(t TestingT, err error) {
	t.Helper()
	if !expect.Success(t, err) {
		t.FailNow()
	}
}
