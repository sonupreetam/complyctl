// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"sort"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/stretchr/testify/require"
)

func TestRenderKeyValuePair(t *testing.T) {

	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "Valid/RenderKeyValuePair",
			key:      "Control ID",
			value:    "r31",
			expected: "Control ID : r31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderKeyValuePair(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetControlRulesColumnsAndRows(t *testing.T) {

	tests := []struct {
		name                 string
		control              Control
		expectedColumnTitles []string
		expectedRows         []table.Row
	}{
		{
			name: "Valid/ControlWithRulesList",
			control: Control{
				ID: "test-control-id",
				Rules: []Rule{
					{ID: "rule-1", Plugin: "Plugin 1"},
					{ID: "rule-2", Plugin: "Plugin 2"},
				},
			},
			expectedColumnTitles: []string{"Rules In Control", "Plugin Used"},
			expectedRows: []table.Row{
				{"rule-1", "Plugin 1"},
				{"rule-2", "Plugin 2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns, rows := getControlRulesColumnsAndRows(tt.control)
			require.Equal(t, len(columns), len(tt.expectedColumnTitles))
			require.Equal(t, len(rows), len(tt.expectedRows))

			sort.Slice(rows, func(i, j int) bool {
				return rows[i][0] < rows[j][0]
			})
			require.Equal(t, rows, tt.expectedRows)
		})
	}
}

func TestGetControlListColumnsAndRows(t *testing.T) {

	tests := []struct {
		name                 string
		controls             []Control
		expectedColumnTitles []string
		expectedRows         []table.Row
	}{
		{
			name: "Valid/ControlsList",
			controls: []Control{
				{
					ID:                  "test-control-id",
					Title:               "Test Control Title",
					ImplemenationStatus: "implemented",
					Rules: []Rule{
						{ID: "rule-1", Plugin: "plugin-1"},
						{ID: "rule-2", Plugin: "plugin-1"},
						{ID: "rule-3", Plugin: "plugin-2"},
					},
				},
			},
			expectedColumnTitles: []string{"Control ID", "Control Title", "Implementation Status", "Plugins Used"},
			expectedRows: []table.Row{
				{"test-control-id", "Test Control Title", "implemented", "plugin-1, plugin-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns, rows := getControlListColumnsAndRows(tt.controls)
			require.Equal(t, len(columns), len(tt.expectedColumnTitles))
			require.Equal(t, len(rows), len(tt.expectedRows))

			sort.Slice(rows, func(i, j int) bool {
				return rows[i][0] < rows[j][0]
			})
			require.Equal(t, rows, tt.expectedRows)
		})
	}
}

func TestGetRuleParametersColumnsAndRows(t *testing.T) {

	tests := []struct {
		name                 string
		ruleInfo             RuleInfo
		setParameters        map[string][]string
		expectedColumnTitles []string
		expectedRows         []table.Row
	}{
		{
			name: "Valid/RuleParametersList",
			ruleInfo: RuleInfo{
				ID:               "rule-1",
				Description:      "Rule 1",
				CheckID:          "check-1",
				CheckDescription: "Check 1",
				Parameters:       []string{"param-1", "param-2"},
			},
			setParameters: map[string][]string{
				"param-1": []string{"param-1-value"},
				"param-2": []string{"param-2-value-1", "param-2-value-2"},
			},
			expectedColumnTitles: []string{"Parameter ID", "Set  Value(s)"},
			expectedRows: []table.Row{
				{"param-1", "param-1-value"},
				{"param-2", "param-2-value-1, param-2-value-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns, rows := getRuleParametersColumnsAndRows(tt.ruleInfo, tt.setParameters)
			require.Equal(t, len(columns), len(tt.expectedColumnTitles))
			require.Equal(t, len(rows), len(tt.expectedRows))

			sort.Slice(rows, func(i, j int) bool {
				return rows[i][0] < rows[j][0]
			})
			require.Equal(t, rows, tt.expectedRows)
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name          string
		inputSlice    []string
		expectedSlice []string
	}{
		{
			name:          "Valid/SliceWithDuplicates",
			inputSlice:    []string{"one", "one", "two", "two"},
			expectedSlice: []string{"one", "two"},
		},
		{
			name:          "Valid/SliceWithDuplicates",
			inputSlice:    []string{},
			expectedSlice: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputSlice := removeDuplicates(tt.inputSlice)
			require.Equal(t, outputSlice, tt.expectedSlice)
		})
	}
}
