// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/cli"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	complytime := cli.New()
	cobra.CheckErr(complytime.ExecuteContext(ctx))
}
