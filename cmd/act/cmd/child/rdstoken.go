package child

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

//Create RDS Token for IAM Authentication
func NewCmdRDSToken() *cobra.Command {
	return builder.NewCmd("rds-token").
		WithDescription("Get RDS Token").
		SetAliases([]string{"rt"}).
		SetFlags().
		RunWithArgsAndCmd(funcGetRDSToken)
}

// Function for rds-token command
func funcGetRDSToken(ctx context.Context, _ io.Writer, cmd *cobra.Command, args []string) error {
	return executor.RunExecutor(ctx, constants.NeedExpiredCheck, func(executor executor.Executor) error {
		var target string
		var err error
		switch len(args) {
		case 0:
			target, err = executor.Runner.ChooseEnv()
			if err != nil {
				return err
			}
		case 1:
			target = args[0]
		default:
			return cmd.Help()
		}

		region := "ap-northeast-2"

		if err := executor.Runner.CopyRDSToken(target, region); err != nil {
			logrus.Errorf(err.Error())
		}
		return nil
	})
}
