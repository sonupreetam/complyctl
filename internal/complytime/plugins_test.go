// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testPluginConfigRoot = filepath.Join("testdata", "complytime", "plugins")

func TestPluginOptions(t *testing.T) {
	tests := []struct {
		name       string
		selections PluginOptions
		wantMap    map[string]string
		wantErr    string
	}{
		{
			name: "Valid/MinimalSelections",
			selections: PluginOptions{
				Workspace: "testworkspace",
				Profile:   "testprofile",
			},
			wantMap: map[string]string{
				"workspace": "testworkspace",
				"profile":   "testprofile",
			},
		},
		{
			name: "Valid/Selections",
			selections: PluginOptions{
				Workspace:      "testworkspace",
				Profile:        "testprofile",
				UserConfigRoot: testPluginConfigRoot,
			},
			wantMap: map[string]string{
				"workspace": "testworkspace",
				"profile":   "testprofile",
				"results":   "results_test.xml",
			},
		},
		{
			name:       "Invalid/MissingOptions",
			selections: PluginOptions{},
			wantErr:    "workspace must be set",
		},
		{
			name: "Invalid/IncorrectOptions",
			selections: PluginOptions{
				Workspace:      "testworkspace",
				Profile:        "testprofile",
				UserConfigRoot: "nonexistpath",
			},
			wantErr: "user config root does not exist",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			err := c.selections.Validate()
			if c.wantErr != "" {
				require.EqualError(t, err, c.wantErr)
			} else {
				require.NoError(t, err)
				gotMap, _ := c.selections.ToMap("openscap")
				require.Equal(t, c.wantMap, gotMap)
			}
		})
	}
}
