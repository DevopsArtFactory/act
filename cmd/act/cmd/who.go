package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

// Assume role with setup
func NewWhoCommand() *cobra.Command {
	return builder.NewCmd("who").
		WithDescription("check the account information of current shell").
		SetFlags().
		RunWithNoArgs(funcWho)
}

// funcWho
func funcWho(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.Who(out)
	})
}
