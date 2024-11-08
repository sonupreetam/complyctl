package cli

import (
	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/spf13/cobra"
	"os"
)

// New creates a new cobra.Command root for ComplyTime
func New() *cobra.Command {
	o := option.Common{}

	o.IOStreams = option.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	cmd := &cobra.Command{
		Use:           "complytime [command]",
		SilenceErrors: false,
		SilenceUsage:  false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	o.BindFlags(cmd.PersistentFlags())
	cmd.AddCommand(versionCmd(&o))

	return cmd
}
