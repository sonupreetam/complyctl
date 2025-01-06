// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/spf13/cobra"

	"path/filepath"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
)

// generateOptions defines options for the generate subcommand
type generateOptions struct {
	*option.Common
	assessmentPlanPath string
}

func setOptsFromArgs(args []string, opts *generateOptions) {
	if len(args) == 1 {
		opts.assessmentPlanPath = filepath.Clean(args[0])
	}
}

// generateCmd creates a new cobra.Command for the generate subcommand
func generateCmd(common *option.Common) *cobra.Command {
	generateOpts := &generateOptions{Common: common}
	return &cobra.Command{
		Use:     "generate",
		Short:   "Generate PVP policy from an assessment plan",
		Example: "complytime generate assessment-plan.json",
		Args:    cobra.RangeArgs(0, 1),
		PreRun: func(cmd *cobra.Command, args []string) {
			setOptsFromArgs(args, generateOpts)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGenerate(cmd, generateOpts)
		},
	}
}

func runGenerate(cmd *cobra.Command, opts *generateOptions) error {

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

	err = manager.GeneratePolicy(cmd.Context(), plugins)
	if err != nil {
		return err
	}
	fmt.Fprintf(opts.Out, "Policy generation completed successfully.")
	return nil
}
