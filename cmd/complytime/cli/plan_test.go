// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"testing"

	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/stretchr/testify/require"

	"github.com/complytime/complytime/cmd/complytime/option"
)

func TestPlansInWorkspace(t *testing.T) {
	// Test that the proper error messages are thrown when
	// working with assessment plans in the user workspace.

	testOpts := &option.ComplyTime{
		UserWorkspace: "doesnotexist",
	}
	wantErr := "error: assessment plan does not exist in workspace doesnotexist: o" +
		"pen doesnotexist/assessment-plan.json: no such file or directory\n\nDid you run the plan command?"
	_, _, gotErr := loadPlan(testOpts, validation.NoopValidator{})
	require.EqualError(t, gotErr, wantErr)

	testOpts.UserWorkspace = "testdata"
	_, gotPath, err := loadPlan(testOpts, validation.NoopValidator{})
	require.NoError(t, err)
	require.Equal(t, "testdata/assessment-plan.json", gotPath)
}

func TestValidatePlan(t *testing.T) {
	tests := []struct {
		name    string
		opts    planOptions
		wantErr string
	}{
		{
			name: "Valid/DefaultOptions",
			opts: planOptions{
				// Set by flag
				output: "-",
			},
		},
		{
			name: "Valid/CorrectPlanOptions",
			opts: planOptions{
				dryRun: true,
				output: "myconfig.yml",
			},
		},
		{
			name: "Invalid/OutNoDryRun",
			opts: planOptions{
				dryRun: false,
				output: "myconfig.yml",
			},
			wantErr: "" +
				"invalid command flags: \"--dry-run\" must be used with \"--out\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.opts)
			err := validatePlan(&tt.opts)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
