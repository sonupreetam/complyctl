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
		wantVersion string
	}{
		{
			name:        "Valid/Default",
			wantVersion: "Version:\tv0.0.0-unknown\n Go Version:\tgo1.22.7\n Git Commit:\t\n Build Date:\t\n",
		},
		{
			name:        "Valid/VariablesSet",
			testVersion: "v0.0.1",
			testCommit:  "commit",
			testDate:    "today",
			testState:   "clean",
			wantVersion: "Version:\tv0.0.1+clean\n Go Version:\tgo1.22.7\n Git Commit:\tcommit\n Build Date:\ttoday\n",
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
			require.Equal(t, c.wantVersion, out.String())

		})
	}
}
