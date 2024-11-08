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
			wantVersion: "Version:\tv0.0.0-unknown\nGo Version:\tgo1.22.7\nGit Commit:\t\nBuild Date:\t\nPlatform:\tlinux/amd64\n",
		},
		{
			name:        "Valid/VariablesSet",
			testVersion: "v0.0.1",
			testCommit:  "commit",
			testDate:    "today",
			testState:   "clean",
			wantVersion: "Version:\tv0.0.1+clean\nGo Version:\tgo1.22.7\nGit Commit:\tcommit\nBuild Date:\ttoday\nPlatform:\tlinux/amd64\n",
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
