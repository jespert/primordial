package expect

import (
	"fmt"
	"strings"
	"testing"
)

var _ TestingT = &testing.T{}

type TestingT interface {
	Helper()
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func Equal[T comparable](t TestingT, want, got T) bool {
	t.Helper()
	if want == got {
		return true
	}

	const msg = "values are not equal"
	report2(t, msg, want, got)
	return false
}

func Success(t TestingT, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
		return false
	}

	return true
}

func Panic(t TestingT, fn func()) (ok bool) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Error(t, "expected panic")
			ok = false
		}
	}()

	ok = true
	fn()
	return ok
}

func report2(t TestingT, msg string, want, got any) {
	t.Errorf("FAIL: %v\n", msg)

	t.Log("Want:")
	if s := fmt.Sprintf("%+v", want); isShortStr(s) {
		t.Logf("   %v\n", s)
	} else {
		t.Logf("\n%v\n\n", s)
	}

	t.Log("Got:")
	if s := fmt.Sprintf("%+v", got); isShortStr(s) {
		t.Logf("    %v\n", s)
	} else {
		t.Logf("\n%v\n", s)
	}
}

func isShortStr(s string) bool {
	// 72 (convenient terminal length) - 8 (prefix) = 64
	const maxShortLen = 64
	return len(s) <= maxShortLen && strings.IndexAny(s, "\n\v\f\r") == -1
}
