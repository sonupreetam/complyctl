// SPDX-License-Identifier: Apache-2.0

package oscap

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConstructScanCommand(t *testing.T) {
	tailoringFilePath := filepath.Join(t.TempDir(), "test-policy.xml")
	tailoringFile, err := os.Create(tailoringFilePath)
	if err != nil {
		t.Fatalf("Failed to create fake tailoring file: %v", err)
	}
	tailoringFile.Close()

	tests := []struct {
		name          string
		openscapFiles map[string]string
		profile       string
		expectedCmd   []string
		expectedErr   bool
	}{
		{
			name: "Valid input with tailoring file",
			openscapFiles: map[string]string{
				"datastream": "test-datastream.xml",
				"policy":     tailoringFilePath,
				"results":    "test-results.xml",
				"arf":        "test-arf.xml",
			},
			profile: "test-profile",
			expectedCmd: []string{
				"oscap",
				"xccdf",
				"eval",
				"--profile",
				"test-profile",
				"--results",
				"test-results.xml",
				"--results-arf",
				"test-arf.xml",
				"--tailoring-file",
				tailoringFilePath,
				"test-datastream.xml",
			},
			expectedErr: false,
		},
		{
			name: "Valid input without tailoring file",
			openscapFiles: map[string]string{
				"datastream": "test-datastream.xml",
				"policy":     "",
				"results":    "test-results.xml",
				"arf":        "test-arf.xml",
			},
			profile: "test-profile",
			expectedCmd: []string{
				"oscap",
				"xccdf",
				"eval",
				"--profile",
				"test-profile",
				"--results",
				"test-results.xml",
				"--results-arf",
				"test-arf.xml",
				"test-datastream.xml",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := constructScanCommand(tt.openscapFiles, tt.profile)
			if (err != nil) != tt.expectedErr {
				t.Errorf("constructScanCommand() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(cmd, tt.expectedCmd) {
				t.Errorf("constructScanCommand() = %v, expected %v", cmd, tt.expectedCmd)
			}
		})
	}
}

// In a more advanced stage we could add tests for the OscapScan function using a minimalistic
// version of a OpenSCAP Datastream, but for now it's not implemented.
