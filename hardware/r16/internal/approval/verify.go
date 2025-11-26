package approval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Verify(t *testing.T, actual string) {
	const testdataDir = "testdata"
	name := t.Name()
	name = strings.TrimPrefix(name, "Test")

	expectedFileName := name + ".expected.txt"
	expectedFilePath := filepath.Join(testdataDir, expectedFileName)
	expected, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read expected file: %v", err)
	}

	if string(expected) != actual {
		t.Errorf("FAIL: Unexpected dump\n\n")
		t.Errorf("Expected:\n%s\nEOF\n\n", expected)
		t.Errorf("Actual:\n%s\nEOF\n\n", actual)
	}
}
