package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

// Initialize act configuration
func NewInitCommand() *cobra.Command {
	return builder.NewCmd("init").
		WithDescription("initialize act command line tool").
		RunWithNoArgs(funcInit)
}

// funcInit
func funcInit(ctx context.Context, _ io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		if err := executor.Runner.InitConfiguration(); err != nil {
			return err
		}

		return nil
	})
}
