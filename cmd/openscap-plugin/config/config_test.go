// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"os/user"
	"path/filepath"
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
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		// Normalizing paths
		{"/foo/bar/../baz", "/foo/baz", false},
		{"./foo/bar", "foo/bar", false},
		{"foo/./bar", "foo/bar", false},
		{"foo/bar/..", "foo", false},
		{"/foo//bar", "/foo/bar", false},
		{"foo//bar//baz", "foo/bar/baz", false},
		{"foo/bar/../../baz", "baz", false},
		{"./../foo", "../foo", false},

		// Expanding paths
		{"~/foo/bar", filepath.Join(homeDir, "foo", "bar"), false},
		{"~", homeDir, false},

		// Weird but valid cases
		{"~weird", "~weird", false}, // not common but possible
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := SanitizePath(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("Expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestSanitizeAndValidatePath tests the SanitizeAndValidatePath function with various
// valid and invalid paths.
func TestSanitizeAndValidatePath(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "testfile")

	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	file.Close()
	defer os.RemoveAll(tempFile)

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

// TestEnsureDirectory tests the ensureDirectory function with various cases.
func TestEnsureDirectory(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		path        string
		expectError bool
	}{
		// Valid cases
		{filepath.Join(tempDir, "absent_dir"), false},   // directory does not exist, should be created
		{filepath.Join(tempDir, "existing_dir"), false}, // directory already exists

		// Invalid cases
		{tempDir + "/invalid\000dir", true}, // invalid directory name
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if tt.path == filepath.Join(tempDir, "existing_dir") {
				// Create directory for existing_dir test
				if err := os.MkdirAll(tt.path, 0750); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
			}

			err := ensureDirectory(tt.path)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			// Check if directory was created
			if !tt.expectError {
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("Expected directory to be created: %s", tt.path)
				}
			}
		})
	}
}

// TestEnsureWorkspace tests the ensureWorkspace function with various cases.
func TestEnsureWorkspace(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		cfg         Config
		expectError bool
	}{
		{
			cfg: Config{
				Files: struct {
					PluginDir  string "yaml:\"plugindir\""
					Workspace  string "yaml:\"workspace\""
					Datastream string "yaml:\"datastream\""
					Results    string "yaml:\"results\""
					ARF        string "yaml:\"arf\""
					Policy     string "yaml:\"policy\""
				}{
					PluginDir: "plugins",
					Workspace: filepath.Join(tempDir, "workspace"),
					Policy:    "policy.yaml",
					Results:   "results.xml",
					ARF:       "arf.xml",
				},
			},
			expectError: false,
		},
		{
			cfg: Config{
				Files: struct {
					PluginDir  string "yaml:\"plugindir\""
					Workspace  string "yaml:\"workspace\""
					Datastream string "yaml:\"datastream\""
					Results    string "yaml:\"results\""
					ARF        string "yaml:\"arf\""
					Policy     string "yaml:\"policy\""
				}{
					PluginDir: "plugins",
					Workspace: filepath.Join(tempDir, "invalid\000workspace"),
					Policy:    "policy.yaml",
					Results:   "results.xml",
					ARF:       "arf.xml",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cfg.Files.Workspace, func(t *testing.T) {
			directories, err := ensureWorkspace(&tt.cfg)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError {
				for _, dir := range directories {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Errorf("Expected directory to be created: %s", dir)
					}
				}
			}
		})
	}
}

// TestDefineFilesPaths tests the defineFilesPaths function with various cases.
func TestDefineFilesPaths(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		cfg         Config
		expectError bool
	}{
		{
			cfg: Config{
				Files: struct {
					PluginDir  string "yaml:\"plugindir\""
					Workspace  string "yaml:\"workspace\""
					Datastream string "yaml:\"datastream\""
					Results    string "yaml:\"results\""
					ARF        string "yaml:\"arf\""
					Policy     string "yaml:\"policy\""
				}{
					PluginDir:  "plugins",
					Workspace:  filepath.Join(tempDir, "workspace"),
					Datastream: filepath.Join(tempDir, "datastream.xml"),
					Results:    "results.xml",
					ARF:        "arf.xml",
					Policy:     "policy.yaml",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cfg.Files.Workspace, func(t *testing.T) {
			_, err := defineFilesPaths(&tt.cfg)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError {
				// Check if the paths are correctly set
				expectedPolicyPath := filepath.Join(tempDir, "workspace", "plugins", "policy", "policy.yaml")
				expectedResultsPath := filepath.Join(tempDir, "workspace", "plugins", "results", "results.xml")
				expectedARFPath := filepath.Join(tempDir, "workspace", "plugins", "results", "arf.xml")

				if tt.cfg.Files.Policy != expectedPolicyPath {
					t.Errorf("Expected policy path: %s, got: %s", expectedPolicyPath, tt.cfg.Files.Policy)
				}
				if tt.cfg.Files.Results != expectedResultsPath {
					t.Errorf("Expected results path: %s, got: %s", expectedResultsPath, tt.cfg.Files.Results)
				}
				if tt.cfg.Files.ARF != expectedARFPath {
					t.Errorf("Expected ARF path: %s, got: %s", expectedARFPath, tt.cfg.Files.ARF)
				}
			}
		})
	}
}

// Tests for ReadConfig are not included because the function relies on other functions
// already tested.
