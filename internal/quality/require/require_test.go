package require_test

import (
	"errors"
	"testing"

	"github.com/jespert/primordial/internal/quality/internal/tmock"
	"github.com/jespert/primordial/internal/quality/require"
)

type T = tmock.T

func TestEqual_ok(t *testing.T) {
	tMock := &T{}
	require.Equal(tMock, true, true)
	if tMock.Terminated {
		t.Error("Expected test to succeed")
	}
}

func TestEqual_fail(t *testing.T) {
	tMock := &T{}
	require.Equal(tMock, true, false)
	if !tMock.Terminated {
		t.Error("Expected test to terminate")
	}
}

func TestSuccess_ok(t *testing.T) {
	tMock := &T{}
	require.Success(tMock, nil)
	if tMock.Terminated {
		t.Error("Expected test to succeed")
	}
}

func TestSuccess_fail(t *testing.T) {
	tMock := &T{}
	require.Success(tMock, errors.New("test error"))
	if !tMock.Terminated {
		t.Error("Expected test to terminate")
	}
}
