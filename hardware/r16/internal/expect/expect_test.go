package expect_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/expect"
)

type result struct {
	ok       bool
	errorOut string
	logOut   string
}

func TestEqual(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		for _, tc := range []struct {
			want, got bool
			result
		}{
			{
				want:   false,
				got:    false,
				result: result{ok: true},
			},
			{
				want: false,
				got:  true,
				result: result{
					ok:       false,
					errorOut: "FAIL: values are not equal\n",
					logOut: "Want:   false\n" +
						"Got:    true\n",
				},
			},
			{
				want: true,
				got:  false,
				result: result{
					ok:       false,
					errorOut: "FAIL: values are not equal\n",
					logOut: "Want:   true\n" +
						"Got:    false\n",
				},
			},
			{
				want: true,
				got:  true,
				result: result{
					ok: true,
				},
			},
		} {
			name := fmt.Sprintf("%v_%v", tc.want, tc.got)
			t.Run(name, func(t *testing.T) {
				tMock := &T{}
				ok := expect.Equal(tMock, tc.want, tc.got)
				validate(t, tMock, tc.want, tc.got, ok, tc.result)
			})
		}
	})
}

// The terminology can be a bit confusing here, because we're talking about
// expectation vs. actual results at two different levels:
// - The user level (want, got)
// - The test validation level (ok, result.ok)
func validate(t *testing.T, tMock *T, want, got any, ok bool, result result) {
	t.Helper()
	if !tMock.HelperCalled {
		t.Errorf("Helper() was not called!")
	}
	if ok != result.ok {
		t.Errorf("Expected that the test returns %v, got %v!", result.ok, ok)
		t.Logf("Want return:  %v\n", result.ok)
		t.Logf("Got return:   %v\n", ok)
	}
	if tMock.Failed && result.ok {
		t.Errorf("Test failed unexpectedly with %v and %v!", want, got)
	}
	if !tMock.Failed && !result.ok {
		t.Errorf("Test succeeded unexpectedly with %v and %v!", want, got)
	}
	if tMock.ErrorOutput.String() != result.errorOut {
		t.Errorf("Unexpected error output")
		t.Errorf("Want:\n%v\n", tMock.ErrorOutput.String())
		t.Errorf("Got:\n%v\n", result.errorOut)
	}
	if tMock.LogOutput.String() != result.logOut {
		t.Errorf("Unexpected log output")
		t.Errorf("Want:\n%v\n", tMock.LogOutput.String())
		t.Errorf("Got:%v\n", result.logOut)
	}
}

var _ expect.TestingT = &T{}

type T struct {
	Failed       bool
	ErrorOutput  bytes.Buffer
	LogOutput    bytes.Buffer
	HelperCalled bool
}

func (t *T) Helper() {
	t.HelperCalled = true
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
