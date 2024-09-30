package cml

import (
	"fmt"
	"strings"
)

// RunFn is a function that runs the command.
//
// Parameters:
//   - p: The program. Assumed to not be nil.
//   - args: The command line arguments.
//
// Returns:
//   - error: Any error that may have occurred.
type RunFn func(p *Program, args []string) error

// Command is a command.
type Command struct {
	// Name is the name of the command.
	Name string

	// Brief is the brief description of the command.
	Brief string

	// Args is the arguments of the command.
	Args *Argument

	// RunFn is the function that runs the command.
	RunFn RunFn
}

// Fix implements the errors.Fixer interface.
func (c *Command) Fix() error {
	if c == nil {
		return nil
	}

	c.Name = strings.TrimSpace(c.Name)
	c.Brief = strings.TrimSpace(c.Brief)

	if c.RunFn == nil {
		c.RunFn = func(_ *Program, _ []string) error {
			return nil
		}
	}

	if c.Args == nil {
		c.Args = NoArguments
	}

	return nil
}

// parse_args is a private function that parses the arguments.
//
// Parameters:
//   - args: The arguments.
//
// Returns:
//   - []string: The parsed arguments.
//   - error: An error if the arguments could not be parsed.
func (c Command) parse_args(args []string) ([]string, error) {
	if c.Args.arg == "" {
		return nil, nil
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("expected argument (%q), got nothing instead", c.Args.arg)
	}

	return args, nil
}
