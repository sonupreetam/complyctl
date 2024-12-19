package config

import (
	"os"
	"testing"
)

// TestSanitizeInput tests the SanitizeInput function with various valid and invalid inputs.
func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		// Valid inputs
		{"valid-input", "valid-input", false},
		{"another_valid.input", "another_valid.input", false},
		{"CAPS_and_numbers123", "CAPS_and_numbers123", false},
		{"mixed-123.UP_case", "mixed-123.UP_case", false},

		// Invalid inputs
		{"invalid/input", "", true},     // contains /
		{"input with spaces", "", true}, // contains spaces
		{"invalid@input", "", true},     // contains @
		{"<invalid>", "", true},         // contains < >
		{";ls", "", true},               // contains ;
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := SanitizeInput(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("Expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestSanitizePath tests the SanitizePath function with various inputs.
func TestSanitizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Normalizing paths
		{"/foo/bar/../baz", "/foo/baz"},
		{"./foo/bar", "foo/bar"},
		{"foo/./bar", "foo/bar"},
		{"foo/bar/..", "foo"},
		{"/foo//bar", "/foo/bar"},
		{"foo//bar//baz", "foo/bar/baz"},
		{"foo/bar/../../baz", "baz"},
		{"./../foo", "../foo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SanitizePath(tt.input)
			if result != tt.expected {
				t.Errorf("Expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestSanitizeAndValidatePath tests the SanitizeAndValidatePath function with various
// valid and invalid paths.
func TestSanitizeAndValidatePath(t *testing.T) {
	tempDir := os.TempDir() + "/test_sanitize_and_validate_path"
	tempFile := tempDir + "/testfile"

	// Setup: create temporary directory and file
	if err := os.MkdirAll(tempDir, 0750); err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	file.Close()
	defer os.RemoveAll(tempDir)

	tests := []struct {
		path        string
		shouldBeDir bool
		expectError bool
		expected    string
	}{
		// Valid cases
		{tempDir, true, false, tempDir},    // directory exists
		{tempFile, false, false, tempFile}, // file exists
		{"/nonexistent", true, true, ""},   // directory does not exist
		{"/nonexistent", false, true, ""},  // file does not exist

		// Invalid cases
		{tempFile, true, true, ""},          // expected directory but found file
		{tempDir, false, true, ""},          // expected file but found directory
		{"/foo/bar/../baz", true, true, ""}, // normalized path does not exist
		{"./foo/bar", true, true, ""},       // relative path does not exist
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result, err := SanitizeAndValidatePath(tt.path, tt.shouldBeDir)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("Expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestEnsureDirectory tests that EnsureDirectory creates a directory if it doesn't exist
// and handles errors correctly.
func TestEnsureDirectory(t *testing.T) {
	tempDir := os.TempDir() + "/test_ensure_directory"

	if _, err := os.Stat(tempDir); err == nil {
		os.RemoveAll(tempDir)
	}

	err := EnsureDirectory(tempDir)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if _, err := os.Stat(tempDir); err != nil {
		t.Errorf("Expected directory to exist, but got error: %v", err)
	}

	defer os.RemoveAll(tempDir)
}

// Tests for ReadConfig, EnsureWorkspace, DefineFilesPaths and ReadConfig functions
// must be created once the configuration file definition is stable.
