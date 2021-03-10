package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

// Assume role with setup
func NewEcrLoginCommand() *cobra.Command {
	return builder.NewCmd("ecr-login").
		WithDescription("login to ECR").
		SetFlags().
		RunWithNoArgs(funcEcrLogin)
}

// funcEcrLogin
func funcEcrLogin(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.EcrLogin(out)
	})
}
