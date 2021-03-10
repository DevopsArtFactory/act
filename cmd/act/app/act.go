package app

import (
	"context"
	"io"

	"github.com/DevopsArtFactory/act/cmd/act/cmd"
)

func Run(out, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	catchCtrlC(cancel)

	c := cmd.NewRootCommand(out, stderr)
	return c.ExecuteContext(ctx)
}
