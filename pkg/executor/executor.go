package executor

import (
	"context"

	"github.com/DevopsArtFactory/act/pkg/builder"
	"github.com/DevopsArtFactory/act/pkg/config"
	"github.com/DevopsArtFactory/act/pkg/runner"
	"github.com/DevopsArtFactory/act/pkg/schema"
)

type Executor struct {
	Runner  runner.Runner
	Context context.Context
}

var NewExecutor = createNewExecutor

// Run executor for command line
func RunExecutor(ctx context.Context, checkKeyExpired bool, action func(Executor) error) error {
	c, err := config.GetConfig()
	if err != nil {
		return err
	}

	awsConfig, err := config.GetCurrentAWSConfig(c.Profile)
	if err != nil {
		return err
	}

	executor, err := createNewExecutor(c)
	if err != nil {
		return err
	}

	if checkKeyExpired {
		if err := executor.Runner.AWSClient.CheckAccessKeyExpired(c.Name, awsConfig.AccessKeyID); err != nil {
			return err
		}
	}

	//Run function with executor
	err = action(*executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

// RunExecutorWithoutCheckingConfig run executor without reading configuration
func RunExecutorWithoutCheckingConfig(ctx context.Context, action func(Executor) error) error {
	executor, err := createNewExecutor(nil)
	if err != nil {
		return err
	}

	//Run function with executor
	err = action(*executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

// RunExecutorReadConfigOnly run executor without checking configuration for command line
func RunExecutorConfigReadOnly(ctx context.Context, action func(Executor) error) error {
	c, err := config.ReadConfigOnly()
	if err != nil {
		return err
	}

	executor, err := createNewExecutor(c)
	if err != nil {
		return err
	}

	//Run function with executor
	err = action(*executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

// Create new executor
func createNewExecutor(config *schema.Config) (*Executor, error) {
	flags, err := builder.ParseFlags()
	if err != nil {
		return nil, err
	}

	executor := Executor{
		Context: context.Background(),
		Runner:  runner.New(flags, config),
	}
	return &executor, nil
}

// alwaysSucceedWhenCancelled makes response true if user canceled
func alwaysSucceedWhenCancelled(ctx context.Context, err error) error {
	// if the context was cancelled act as if all is well
	if err != nil && ctx.Err() == context.Canceled {
		return nil
	}
	return err
}
