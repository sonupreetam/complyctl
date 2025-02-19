// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
	"github.com/complytime/complytime/internal/terminal"
)

// listCmd creates a new cobra.Command for the "list" subcommand
func listCmd(common *option.Common, logger hclog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list [flags]",
		Short:        "List information about supported frameworks and components.",
		SilenceUsage: true,
		Example:      "complytime list",
		Args:         cobra.NoArgs,
		PreRun:       func(_ *cobra.Command, _ []string) { enableDebug(logger, common) },
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := runList(common, logger); err != nil {
				logger.Error(err.Error())
			}
			return nil
		},
	}
	return cmd
}

func runList(opts *option.Common, logger hclog.Logger) error {
	appDir, err := complytime.NewApplicationDirectory(true)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("using application directory: %s", appDir.AppDir()))

	frameworks, err := complytime.LoadFrameworks(appDir)
	if err != nil {
		return err
	}

	model := terminal.ShowDefinitionTable(frameworks)
	if _, err := tea.NewProgram(model, tea.WithOutput(opts.Out)).Run(); err != nil {
		return fmt.Errorf("failed to display component list: %w", err)
	}
	return nil
}
