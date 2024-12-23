// SPDX-License-Identifier: Apache-2.0

package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteVersion(t *testing.T) {
	tests := []struct {
		name        string
		testVersion string
		testCommit  string
		testDate    string
		testState   string
		assertFunc  func(string) bool
	}{
		{
			name: "Valid/Default",
			assertFunc: func(s string) bool {
				return strings.Contains(s, "Version:\tv0.0.0-unknown\n") &&
					strings.Contains(s, "Git Commit:\t\nBuild Date:\t\n")
			},
		},
		{
			name:        "Valid/VariablesSet",
			testVersion: "v0.0.1",
			testCommit:  "commit",
			testDate:    "today",
			testState:   "clean",
			assertFunc: func(s string) bool {
				return strings.Contains(s, "v0.0.1+clean") &&
					strings.Contains(s, "Git Commit:\tcommit\n") &&
					strings.Contains(s, "Build Date:\ttoday\n")
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			out := new(strings.Builder)
			if c.testVersion != "" {
				version = c.testVersion
			}
			buildDate = c.testDate
			commit = c.testCommit
			gitTreeState = c.testState
			err := WriteVersion(out)

			require.NoError(t, err)
			require.True(t, c.assertFunc(out.String()))
		})
	}
}
