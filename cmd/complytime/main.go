package main

import (
	"context"
	"github.com/complytime/complytime/cmd/complytime/cli"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	complytime := cli.New()
	cobra.CheckErr(complytime.ExecuteContext(ctx))
}
