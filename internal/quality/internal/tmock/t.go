package tmock

import (
	"bytes"
	"fmt"
)

type T struct {
	Failed       bool
	ErrorOutput  bytes.Buffer
	LogOutput    bytes.Buffer
	HelperCalled bool
	Terminated   bool
}

func (t *T) Helper() {
	t.HelperCalled = true
}

func (t *T) Error(args ...interface{}) {
	t.Failed = true
	for _, arg := range args {
		// There's nothing we can do about IO errors here.
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
		// There's nothing we can do about IO errors here.
		_, _ = fmt.Fprintf(&t.LogOutput, "%v", arg)
	}
}

func (t *T) Logf(format string, args ...interface{}) {
	// There's nothing we can do about IO errors here.
	_, _ = fmt.Fprintf(&t.LogOutput, format, args...)
}

func (t *T) FailNow() {
	t.Terminated = true
}
