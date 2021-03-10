package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/version"
)

// Get act version
func NewVersionCommand() *cobra.Command {
	return builder.NewCmd("version").
		WithDescription("Print the version information").
		SetAliases([]string{"v"}).
		RunWithNoArgs(funcVersion)
}

// funcVersion
func funcVersion(_ context.Context, _ io.Writer) error {
	return version.Controller{}.Print(version.Get())
}
