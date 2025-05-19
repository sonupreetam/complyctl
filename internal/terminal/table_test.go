// SPDX-License-Identifier: Apache-2.0

package terminal

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/stretchr/testify/require"
)

func TestModelView(t *testing.T) {

	tests := []struct {
		name       string
		model      Model
		wantOutput string
	}{
		{
			name: "Valid/Table",
			model: Model{Table: table.New(
				table.WithColumns([]table.Column{{Title: "Column A", Width: 10}}),
				table.WithRows([]table.Row{{"Row A"}}),
				table.WithHeight(7),
			)},
			wantOutput: testTable,
		},
		{
			name: "Valid/TableWithHelp",
			model: Model{
				Table: table.New(
					table.WithColumns([]table.Column{{Title: "Column A", Width: 10}}),
					table.WithRows([]table.Row{{"Row A"}}),
					table.WithHeight(7)),
				HelpMsg: "help message here",
			},
			wantOutput: testTableWithHelp,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			gotOutput := c.model.View()
			require.Equal(t, c.wantOutput, gotOutput)
		})
	}
}

var (
	testTable = `┌────────────┐
│ Column A   │
│ Row A      │
│            │
│            │
│            │
│            │
│            │
└────────────┘
`

	testTableWithHelp = `┌────────────┐
│ Column A   │
│ Row A      │
│            │
│            │
│            │
│            │
│            │
└────────────┘
help message here
`
)
