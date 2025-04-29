// SPDX-License-Identifier: Apache-2.0

package oscap

import (
	"reflect"
	"testing"
)

func TestConstructScanCommand(t *testing.T) {
	tests := []struct {
		name          string
		openscapFiles map[string]string
		profile       string
		expectedCmd   []string
	}{
		{
			name: "Scan command contruction",
			openscapFiles: map[string]string{
				"datastream": "test-datastream.xml",
				"policy":     "test-policy.xml",
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
				"test-policy.xml",
				"test-datastream.xml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := constructScanCommand(tt.openscapFiles, tt.profile)
			if !reflect.DeepEqual(cmd, tt.expectedCmd) {
				t.Errorf("constructScanCommand() = %v, expected %v", cmd, tt.expectedCmd)
			}
		})
	}
}

// In a more advanced stage we could add tests for the OscapScan function using a minimalistic
// version of a OpenSCAP Datastream, but for now it's not implemented.

func TestConstructGenerateFixCommand(t *testing.T) {
	tests := []struct {
		name          string
		fixType       string
		output        string
		profile       string
		tailoringFile string
		datastream    string
		expectedCmd   []string
	}{
		{
			name:          "Genereate fix command construction",
			fixType:       "bash",
			output:        "test-remediation-script.sh",
			profile:       "test-profile",
			tailoringFile: "test-policy.xml",
			datastream:    "test-datastream.xml",
			expectedCmd: []string{
				"oscap",
				"xccdf",
				"generate",
				"fix",
				"--fix-type", "bash",
				"--output", "test-remediation-script.sh",
				"--profile", "test-profile",
				"--tailoring-file", "test-policy.xml",
				"test-datastream.xml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := constructGenerateFixCommand(tt.fixType, tt.output, tt.profile, tt.tailoringFile, tt.datastream)
			if !reflect.DeepEqual(cmd, tt.expectedCmd) {
				t.Errorf("constructGenerateFixCommand() = %v, expected %v", cmd, tt.expectedCmd)
			}
		})
	}
}
