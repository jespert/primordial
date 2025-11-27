package approval

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/assert"
)

var _ TestingT = &testing.T{}

type TestingT interface {
	Cleanup(func())
	Helper()
	Name() string
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

type Verifier struct {
	t                TestingT
	expectedFilePath string
	actualFilePath   string
	expectedFile     *os.File
	actualFile       *os.File
	verifyCalled     bool
}

func NewTextVerifier(t TestingT) *Verifier {
	t.Helper()

	// Name the files after the test function but remove the "Test" prefix,
	// which only adds noise.
	name := t.Name()
	name = strings.TrimPrefix(name, "Test")

	verifier := &Verifier{
		t:                t,
		expectedFilePath: filepath.Join(testdataDir, name+".expected.txt"),
		actualFilePath:   filepath.Join(testdataDir, name+".actual.txt"),
	}

	// Install the clean-up function eagerly to cover all failure modes.
	t.Cleanup(verifier.cleanup())

	var err error
	verifier.expectedFile, err = os.Open(verifier.expectedFilePath)
	if err != nil {
		t.Fatal("Failed to read expected data:", err)
	}

	verifier.actualFile, err = os.Create(verifier.actualFilePath)
	if err != nil {
		t.Fatal("Failed to write actual data:", err)
	}

	return verifier
}

func (v *Verifier) Writer() io.Writer {
	// Don't grant the user the ability to close the writer.
	return v.actualFile
}

func (v *Verifier) Verify() bool {
	t := v.t
	t.Helper()

	v.verifyCalled = true

	// Comparing file sizes for equality first is inexpensive and
	// covers a lot of cases in practice.
	if v.filesHaveSameSize() && v.filesHaveSameContent() {
		// The files are identical. Close them.
		if err := v.actualFile.Close(); err != nil {
			t.Fatal("Failed to close actual file:", err)
		}
		v.actualFile = nil

		if err := v.expectedFile.Close(); err != nil {
			t.Fatal("Failed to close expected file:", err)
		}
		v.expectedFile = nil

		// Delete the actual file to reduce clutter.
		err := os.Remove(v.actualFilePath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal("Failed to delete actual file:", err)
		}

		return true
	}

	// The files are different.
	t.Error("FAIL: Unexpected result.\n")
	t.Logf("Expected: %v\n", v.expectedFilePath)
	t.Logf("Actual: %v\n", v.actualFilePath)
	return false
}

func (v *Verifier) filesHaveSameSize() bool {
	t := v.t

	actualStat, err := v.actualFile.Stat()
	if err != nil {
		t.Fatal("Failed to stat actual data:", err)
		return false // unreachable
	}

	expectedStat, err := v.expectedFile.Stat()
	if err != nil {
		t.Fatal("Failed to stat actual data:", err)
		return false // unreachable
	}

	return actualStat.Size() == expectedStat.Size()
}

func (v *Verifier) filesHaveSameContent() bool {
	t := v.t

	_, err := v.actualFile.Seek(0, io.SeekStart)
	if err != nil {
		return false
	}

	const bufferSize = 128 * 1024
	buf1 := make([]byte, bufferSize)
	buf2 := make([]byte, bufferSize)

	for {
		n1, err1 := v.expectedFile.Read(buf1)
		n2, err2 := v.actualFile.Read(buf2)
		if n1 != n2 {
			return false
		}
		if bytes.Compare(buf1[:n1], buf2[:n2]) != 0 {
			return false
		}
		if err1 == io.EOF && err2 == io.EOF {
			return true
		}
		if err1 != nil && err1 != io.EOF {
			t.Error("Failed to read expected data:", err)
			return false
		}
		if err2 != nil && err2 != io.EOF {
			t.Error("Failed to read actual data:", err)
			return false
		}
	}
}

func (v *Verifier) cleanup() func() {
	return func() {
		t := v.t
		if v.expectedFile != nil {
			err := v.expectedFile.Close()
			if err != nil {
				t.Error("Failed to close expected file:", err)
			}
		}
		if v.actualFile != nil {
			err := v.actualFile.Close()
			if err != nil {
				t.Error("Failed to close actual file:", err)
			}
		}
		if !v.verifyCalled {
			v.t.Error("FAIL: Test did not call Verify().")
		}
	}
}

func init() {
	cwd, err := os.Getwd()
	assert.Success(err)

	testdataDir = filepath.Join(cwd, "testdata")
}

var testdataDir string
