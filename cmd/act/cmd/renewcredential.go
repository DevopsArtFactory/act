package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

// renew credentials
func NewRenewCredentialsCommand() *cobra.Command {
	return builder.NewCmd("renew-credential").
		WithDescription("recreates aws credential of profile").
		SetFlags().
		RunWithNoArgs(funcRenewCredentials)
}

// funcRenewCredentials
func funcRenewCredentials(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, constants.SkipExpiredCheck, func(executor executor.Executor) error {
		return executor.Runner.RenewCredentials(out)
	})
}
