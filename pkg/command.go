package pkg

import "strings"

// RunFunc is a function that will be executed when the command is called.
//
// Parameters:
//   - p: The program that the command is being executed on.
//   - args: The arguments that were passed to the command.
//
// Returns:
//   - error: An error if the command failed to execute.
type RunFunc func(p *Program, args []string) error

type Command struct {
	// Name is the name of the command.
	Name string

	// Usage is the usage of the command.
	Usage string

	// Brief is a brief description of the command.
	Brief string

	// Description is a description of the command.
	Description []string

	// Argument is the argument of the command.
	Argument *Argument

	// Run is the function that will be executed when the command is called.
	Run RunFunc
}

// fix fixes the command by trimming all the strings and setting default values.
func (c *Command) fix() {
	c.Name = strings.TrimSpace(c.Name)
	c.Usage = strings.TrimSpace(c.Usage)
	c.Brief = strings.TrimSpace(c.Brief)

	if c.Argument == nil {
		c.Argument = NoArgument
	}

	if c.Run == nil {
		c.Run = func(p *Program, args []string) error {
			return nil
		}
	}
}
