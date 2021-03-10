package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/child"
)

// Get information with act
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get token or information with act",
	}

	cmd.AddCommand(child.NewCmdRDSToken())
	return cmd
}
