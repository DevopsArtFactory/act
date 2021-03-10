package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

// Assume role with setup
func NewSetupCommand() *cobra.Command {
	return builder.NewCmd("setup").
		WithDescription("create assume credentials for multi-account").
		SetFlags().
		RunWithArgsAndCmd(funcSetup)
}

// funcSetup
func funcSetup(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return cmd.Help()
	}

	return executor.RunExecutor(ctx, constants.NeedExpiredCheck, func(executor executor.Executor) error {
		return executor.Runner.Setup(out, args)
	})
}
