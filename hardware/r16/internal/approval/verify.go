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
		t.Fatalf("Failed to read expected data: %v", err)
	}

	actualFileName := name + ".actual.txt"
	actualFilePath := filepath.Join(testdataDir, actualFileName)
	if string(expected) == actual {
		if err := os.Remove(actualFilePath); err != nil {
			t.Logf("Failed to remove actual data: %v", err)
		}

		// Test passed.
		return
	}

	if err := os.WriteFile(
		actualFilePath,
		[]byte(actual),
		0644,
	); err != nil {
		t.Fatalf("Failed to write actual data: %v", err)
	}

	// Reporting absolute paths improves the developer UX.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"Failed to get current working directory: %v",
			err,
		)
	}

	expectedFilePath = filepath.Join(cwd, expectedFilePath)
	actualFilePath = filepath.Join(cwd, actualFilePath)

	t.Errorf("FAIL: Unexpected result\n\n")
	t.Logf("Expected: %v\n", expectedFilePath)
	t.Logf("Actual: %v\n", actualFilePath)
}
