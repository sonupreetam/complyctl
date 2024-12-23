// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
)

// scanOptions defined options for the scan subcommand.
type scanOptions struct {
	*option.Common
	assessmentPlanPath string
}

// scanCmd creates a new cobra.Command for the version subcommand.
func scanCmd(common *option.Common) *cobra.Command {
	scanOpts := &scanOptions{Common: common}
	return &cobra.Command{
		Use:          "scan [flags] <assessment_plan_path>",
		Short:        "Scan environment with assessment plan",
		Example:      "complytime scan assessment-plan.json",
		SilenceUsage: true,
		Args:         cobra.RangeArgs(0, 1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				scanOpts.assessmentPlanPath = filepath.Clean(args[0])
			}
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runScan(cmd, scanOpts)
		},
	}
}

func runScan(cmd *cobra.Command, opts *scanOptions) error {
	// Adding this message to the user for now because assessment Plans are unused
	if opts.assessmentPlanPath != "" {
		_, _ = fmt.Fprintf(opts.Out, "OSCAL Assessment Plans are not supported yet...\nThe file %s will not be used.\n", opts.assessmentPlanPath)
	}

	// Create the application directory if it does not exist
	appDir, err := complytime.NewApplicationDirectory(true)
	if err != nil {
		return err
	}

	cfg, err := complytime.Config(appDir)
	if err != nil {
		return err
	}
	manager, err := framework.NewPluginManager(cfg)
	if err != nil {
		return fmt.Errorf("error initializing plugin manager: %w", err)
	}
	plugins, err := manager.LaunchPolicyPlugins()
	if err != nil {
		return err
	}
	// Ensure all the plugins launch above are cleaned up
	defer manager.Clean()

	allResults, err := manager.AggregateResults(cmd.Context(), plugins)
	if err != nil {
		return err
	}

	// This is a temporary solution for results processing.
	return displayResults(opts.Out, allResults)
}

// displayResults write the results from the scan with the given io.Writer.
func displayResults(writer io.Writer, allResults []policy.PVPResult) error {
	_, _ = fmt.Fprintf(writer, "Processing %d result(s)...\n", len(allResults))
	for _, results := range allResults {
		for _, observation := range results.ObservationsByCheck {
			_, _ = fmt.Fprintf(writer, "Observation: %s\n", observation.Title)
			for _, sub := range observation.Subjects {
				_, _ = fmt.Fprintf(writer, "Subject: %s\n", sub.Title)
				_, _ = fmt.Fprintf(writer, "Resource: %s\n", sub.ResourceID)
				_, _ = fmt.Fprintf(writer, "Result: %s\n", sub.Result.String())
				_, _ = fmt.Fprintf(writer, "Reason: %s\n\n", sub.Reason)
			}
		}
		_, _ = fmt.Fprintf(writer, "\n")
	}
	return nil
}
