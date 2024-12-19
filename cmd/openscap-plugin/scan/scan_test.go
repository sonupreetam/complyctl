// SPDX-License-Identifier: Apache-2.0
package scan

import (
	"fmt"
	"os"
	"testing"
)

func TestIsXMLFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		want      bool
		expectErr bool
	}{
		{
			name:      "Valid XML file",
			filePath:  "testdata/valid.xml",
			want:      true,
			expectErr: false,
		},
		{
			name:      "Invalid XML file",
			filePath:  "testdata/invalid.xml",
			want:      false,
			expectErr: true,
		},
		{
			name:      "Non-existent file",
			filePath:  "testdata/nonexistent.xml",
			want:      false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isXMLFile(tt.filePath)
			if (err != nil) != tt.expectErr {
				t.Errorf("isXMLFile() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("isXMLFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDataStream(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		setup     func()
		want      string
		expectErr bool
	}{
		{
			name:     "Valid datastream file",
			filePath: "testdata/valid.xml",
			setup: func() {
				if err := os.WriteFile("testdata/valid.xml", []byte(`<root></root>`), 0600); err != nil {
					t.Fatalf("Failed to write valid.xml: %v", err)
				}
			},
			want:      "testdata/valid.xml",
			expectErr: false,
		},
		{
			name:     "Invalid datastream file",
			filePath: "testdata/invalid.xml",
			setup: func() {
				if err := os.WriteFile("testdata/invalid.xml", []byte(`<root>`), 0600); err != nil {
					t.Fatalf("Failed to write invalid.xml: %v", err)
				}
			},
			want:      "",
			expectErr: true,
		},
		{
			name:      "Non-existent datastream file",
			filePath:  "testdata/nonexistent.xml",
			setup:     func() {},
			want:      "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := validateDataStream(tt.filePath)
			if (err != nil) != tt.expectErr {
				t.Errorf("validateDataStream() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateDataStream() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestValidateTailoringFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		setup     func()
		want      string
		expectErr bool
	}{
		{
			name:     "Valid tailoring file",
			filePath: "testdata/valid.xml",
			setup: func() {
				if err := os.WriteFile("testdata/valid.xml", []byte(`<root></root>`), 0600); err != nil {
					t.Fatalf("Failed to write valid.xml: %v", err)
				}
			},
			want:      "testdata/valid.xml",
			expectErr: false,
		},
		{
			name:     "Invalid tailoring file",
			filePath: "testdata/invalid.xml",
			setup: func() {
				if err := os.WriteFile("testdata/invalid.xml", []byte(`<root>`), 0600); err != nil {
					t.Fatalf("Failed to write invalid.xml: %v", err)
				}
			},
			want:      "",
			expectErr: true,
		},
		{
			name:      "Non-existent tailoring file",
			filePath:  "testdata/nonexistent.xml",
			setup:     func() {},
			want:      "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := validateTailoringFile(tt.filePath)
			if (err != nil) != tt.expectErr {
				t.Errorf("validateTailoringFile() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateTailoringFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupTestFiles() error {
	if err := os.MkdirAll("testdata", os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile("testdata/valid.xml", []byte(`<root></root>`), 0600); err != nil {
		return err
	}
	if err := os.WriteFile("testdata/invalid.xml", []byte(`<root>`), 0600); err != nil {
		return err
	}
	return nil
}

func teardownTestFiles() {
	os.RemoveAll("testdata")
}

func TestMain(m *testing.M) {
	if err := setupTestFiles(); err != nil {
		fmt.Printf("Failed to setup test files: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	teardownTestFiles()
	os.Exit(code)
}

// ScanSystem function is not tested because it is a high-level function that uses other functions
// already tested above or in other packages.
