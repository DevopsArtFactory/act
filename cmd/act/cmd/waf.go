package cmd

import (
	"context"
	"errors"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/act/cmd/act/cmd/builder"
	"github.com/DevopsArtFactory/act/pkg/executor"
)

// Describe detailed information of web acl
func NewCmdDescribeWebACL() *cobra.Command {
	return builder.NewCmd("describe-web-acl").
		WithDescription("retrieve detailed list of web acl").
		SetAliases([]string{"dwa"}).
		SetFlags().
		RunWithArgs(funcDescribeWebACL)
}

// Function for describe-web-acl command
func funcDescribeWebACL(ctx context.Context, out io.Writer, args []string) error {
	if len(args) > 1 {
		return errors.New("usage: act describe-web-acl(dwa) [WAF name]")
	}
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		if err := executor.Runner.DescribeWebACL(out, args); err != nil {
			logrus.Errorf(err.Error())
		}
		return nil
	})
}

// Describe detailed information of web acl
func NewCmdHasIP() *cobra.Command {
	return builder.NewCmd("has-ip").
		WithDescription("check if ip is registered in the web acl").
		SetFlags().
		RunWithArgs(funcHasIP)
}

// Function for has-ip command
func funcHasIP(ctx context.Context, out io.Writer, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: act has-ip [IP1] [IP2]... ")
	}
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		if err := executor.Runner.HasIP(out, args); err != nil {
			logrus.Errorf(err.Error())
		}
		return nil
	})
}
