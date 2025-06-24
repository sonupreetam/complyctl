// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/spf13/cobra"

	"github.com/complytime/complyctl/cmd/complyctl/option"

	"github.com/complytime/complyctl/internal/version"
)

// versionCmd creates a new cobra.Command for the version subcommand.
func versionCmd(common *option.Common) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return version.WriteVersion(common.Out)
		},
	}
}
