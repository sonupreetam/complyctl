// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"os"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/pkg/log"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

func enableDebug(logger hclog.Logger, opts *option.Common) {
	if opts.Debug {
		logger.SetLevel(hclog.Debug)
	}
}

// New creates a new cobra.Command root for ComplyTime
func New() *cobra.Command {

	logger := log.NewLogger(os.Stdout)

	cmd := &cobra.Command{
		Use:           "complytime [command]",
		SilenceErrors: true,
		SilenceUsage:  false,
	}

	cmd.Context()

	opts := option.Common{
		Output: option.Output{
			Out:    cmd.OutOrStdout(),
			ErrOut: cmd.ErrOrStderr(),
		},
	}
	opts.BindFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		versionCmd(&opts),
		scanCmd(&opts, logger),
		generateCmd(&opts, logger),
		planCmd(&opts, logger),
		listCmd(&opts, logger),
	)

	return cmd
}
