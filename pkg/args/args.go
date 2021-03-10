package args

import (
	"fmt"
)

type Argument struct {
	Profile string
	Command string
	Args    []string
}

// Parse parse os arguments to Argument struct
func Parse(args []string) (*Argument, error) {
	a := Argument{}
	if !IsValid(args) {
		return nil, fmt.Errorf("usage: act exec [profile] -- [command]")
	}

	parseArgument(args, &a)

	return &a, nil
}

// IsValid checks if args is valid or not
func IsValid(args []string) bool {
	return len(args) >= 3
}

// parseArgument parsing arguments
func parseArgument(args []string, a *Argument) {
	a.Profile = args[0]
	a.Command = args[1]
	a.Args = args[2:]
}
