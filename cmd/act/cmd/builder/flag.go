package builder

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/tools"
)

type Flag struct {
	Name               string
	Shorthand          string
	Usage              string
	Value              interface{}
	DefValue           interface{}
	DefValuePerCommand map[string]interface{}
	FlagAddMethod      string
	DefinedOn          []string
	Hidden             bool

	pflag *pflag.Flag
}

// FlagRegistry is a list of all act CLI flags.
var FlagRegistry = []Flag{
	{
		Name:          "region",
		Shorthand:     "r",
		Usage:         "Run command to specific region",
		Value:         aws.String(constants.EmptyString),
		DefValue:      "ap-northeast-2",
		FlagAddMethod: "StringVar",
		DefinedOn:     []string{"rds-token", "who", "describe-web-acl"},
	},
	{
		Name:          "duration",
		Shorthand:     "d",
		Usage:         "Duration of IAM assume role session time",
		Value:         aws.Int(constants.DefaultDuration),
		DefValue:      constants.DefaultDuration,
		FlagAddMethod: "IntVar",
		DefinedOn:     []string{"setup"},
	},
	{
		Name:          "profile",
		Shorthand:     "p",
		Usage:         "Profile you want to use",
		Value:         aws.String("default"),
		DefValue:      "default",
		FlagAddMethod: "StringVar",
		DefinedOn:     []string{"rds-token", "setup", "list", "renew-credential", "describe-web-acl", "has-ip", "start", "stop", "status", "exec", "ecr-login"},
	},
	{
		Name:          "raw-output",
		Usage:         "show raw output without copying to clipboard",
		Value:         aws.Bool(false),
		DefValue:      false,
		FlagAddMethod: "BoolVar",
		DefinedOn:     []string{"setup"},
	},
}

func (fl *Flag) flag() *pflag.Flag {
	if fl.pflag != nil {
		return fl.pflag
	}

	inputs := []interface{}{fl.Value, fl.Name}
	if fl.FlagAddMethod != "Var" {
		inputs = append(inputs, fl.DefValue)
	}
	inputs = append(inputs, fl.Usage)

	fs := pflag.NewFlagSet(fl.Name, pflag.ContinueOnError)
	reflect.ValueOf(fs).MethodByName(fl.FlagAddMethod).Call(reflectValueOf(inputs))
	f := fs.Lookup(fl.Name)
	f.Shorthand = fl.Shorthand
	f.Hidden = fl.Hidden

	fl.pflag = f
	return f
}

func reflectValueOf(values []interface{}) []reflect.Value {
	var results []reflect.Value
	for _, v := range values {
		results = append(results, reflect.ValueOf(v))
	}
	return results
}

//Add command flags
func SetCommandFlags(cmd *cobra.Command) {
	var flagsForCommand []*Flag
	for i := range FlagRegistry {
		fl := &FlagRegistry[i]

		if tools.IsStringInArray(cmd.Use, fl.DefinedOn) {
			cmd.Flags().AddFlag(fl.flag())
			flagsForCommand = append(flagsForCommand, fl)
		}
	}

	// Apply command-specific default values to flags.
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Update default values.
		for _, fl := range flagsForCommand {
			viper.BindPFlag(fl.Name, cmd.Flags().Lookup(fl.Name))
		}

		// Since PersistentPreRunE replaces the parent's PersistentPreRunE,
		// make sure we call it, if it is set.
		if parent := cmd.Parent(); parent != nil {
			if preRun := parent.PersistentPreRunE; preRun != nil {
				if err := preRun(cmd, args); err != nil {
					return err
				}
			} else if preRun := parent.PersistentPreRun; preRun != nil {
				preRun(cmd, args)
			}
		}

		return nil
	}
}
