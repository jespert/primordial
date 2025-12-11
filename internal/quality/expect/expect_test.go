package expect_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jespert/primordial/internal/quality/expect"
)

type result struct {
	ok       bool
	errorOut string
	logOut   string
}

func TestEqual_bool(t *testing.T) {
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
}

func TestEqual_string(t *testing.T) {
	for _, tc := range []struct {
		name      string
		want, got string
		result
	}{
		{
			name:   "same short string",
			want:   "foo",
			got:    "foo",
			result: result{ok: true},
		},
		{
			name: "different short string",
			want: "foo",
			got:  "bar",
			result: result{
				ok:       false,
				errorOut: "FAIL: values are not equal\n",
				logOut: "Want:   foo\n" +
					"Got:    bar\n",
			},
		},
		{
			name:   "same long string",
			want:   longStr1,
			got:    longStr1,
			result: result{ok: true},
		},
		{
			name: "short and long string",
			want: "foo",
			got:  longStr2,
			result: result{
				ok:       false,
				errorOut: "FAIL: values are not equal\n",
				logOut: "Want:   foo\n" +
					"Got:\n" +
					"Donuts in the break room pulling teeth,\n" +
					"nor strategic staircase,\n" +
					"yet high touch client.\n",
			},
		},
		{
			name: "long and short string",
			want: longStr1,
			got:  "foo",
			result: result{
				ok:       false,
				errorOut: "FAIL: values are not equal\n",
				logOut: "Want:\n" +
					"Lorem ipsum dolor sit amet,\n" +
					"consectetur adipiscing elit.\n" +
					"Donec auctor, lorem quis tincidunt consequat,\n" +
					"elit elit dignissim elit.\n" +
					"\n" +
					"Got:    foo\n",
			},
		},
		{
			name: "different long strings",
			want: longStr1,
			got:  longStr2,
			result: result{
				ok:       false,
				errorOut: "FAIL: values are not equal\n",
				logOut: "Want:\n" +
					"Lorem ipsum dolor sit amet,\n" +
					"consectetur adipiscing elit.\n" +
					"Donec auctor, lorem quis tincidunt consequat,\n" +
					"elit elit dignissim elit.\n" +
					"\n" +
					"Got:\n" +
					"Donuts in the break room pulling teeth,\n" +
					"nor strategic staircase,\n" +
					"yet high touch client.\n",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tMock := &T{}
			ok := expect.Equal(tMock, tc.want, tc.got)
			validate(t, tMock, tc.want, tc.got, ok, tc.result)
		})
	}
}

func TestPanic_success(t *testing.T) {
	tMock := &T{}
	ok := expect.Panic(tMock, func() { panic("test panic") })
	if !ok {
		t.Error("Expected panic")
	}
}

func TestPanic_fail(t *testing.T) {
	tMock := &T{}
	ok := expect.Panic(tMock, func() {})
	if ok {
		t.Error("Unexpected panic")
	}
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

const longStr1 = `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.
Donec auctor, lorem quis tincidunt consequat,
elit elit dignissim elit.`

const longStr2 = `Donuts in the break room pulling teeth,
nor strategic staircase,
yet high touch client.`
