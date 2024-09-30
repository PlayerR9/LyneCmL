package cml

import "strings"

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

	return nil
}
