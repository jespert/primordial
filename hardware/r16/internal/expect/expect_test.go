package expect_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/expect"
)

type result struct {
	expected bool
	errorOut string
	logOut   string
}

func TestEqual(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		for _, tc := range []struct {
			x, y bool
			result
		}{
			{
				x:      false,
				y:      false,
				result: result{expected: true},
			},
			{
				x: false,
				y: true,
				result: result{
					expected: false,
					errorOut: "FAIL: values are not equal\n",
					logOut: "Want:   false\n" +
						"Got:    true\n",
				},
			},
			{
				x: true,
				y: false,
				result: result{
					expected: false,
					errorOut: "FAIL: values are not equal\n",
					logOut: "Want:   true\n" +
						"Got:    false\n",
				},
			},
			{
				x: true,
				y: true,
				result: result{
					expected: true,
				},
			},
		} {
			name := fmt.Sprintf("%v_%v", tc.x, tc.y)
			t.Run(name, func(t *testing.T) {
				tMock := &T{}
				actual := expect.Equal(tMock, tc.x, tc.y)
				report(tMock, tc.x, tc.y, actual, tc.result)
			})
		}
	})
}

func report(t *T, x, y any, ok bool, result result) {
	if ok != result.expected {
		t.Errorf("Expected %v, got %v!", result.expected, ok)
		t.Logf("Want:\n%v\n", result.expected)
		t.Logf("Got:\n%v\n", ok)
	}
	if t.Failed == result.expected {
		t.Errorf("Failed with %v and %v!", x, y)
	}
	if t.ErrorOutput.String() != result.errorOut {
		t.Errorf("Unexpected error output")
		t.Errorf("Want:\n%v\n", t.ErrorOutput.String())
		t.Errorf("Got:\n%v\n", result.errorOut)
	}
	if t.LogOutput.String() != result.logOut {
		t.Errorf("Unexpected log output")
		t.Errorf("Want:\n%v\n", t.LogOutput.String())
		t.Errorf("Got:%v\n", result.logOut)
	}
}

var _ expect.TestingT = &T{}

type T struct {
	Failed      bool
	ErrorOutput bytes.Buffer
	LogOutput   bytes.Buffer
}

func (t *T) Error(args ...interface{}) {
	for _, arg := range args {
		_, _ = fmt.Fprintf(&t.ErrorOutput, "%v", arg)
	}
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.Failed = true

	// There's nothing we can do about IO errors here.
	_, _ = fmt.Fprintf(&t.ErrorOutput, format, args...)
}

func (t *T) Log(args ...interface{}) {
	for _, arg := range args {
		_, _ = fmt.Fprintf(&t.LogOutput, "%v", arg)
	}
}

func (t *T) Logf(format string, args ...interface{}) {
	// There's nothing we can do about IO errors here.
	_, _ = fmt.Fprintf(&t.LogOutput, format, args...)
}
