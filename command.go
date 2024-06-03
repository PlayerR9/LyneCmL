package LyneCml

// CommandRunFunc is a function that runs a command.
//
// Parameters:
//   - p: The program that the command is running on.
//   - args: The arguments that the command is running with.
//
// Returns:
//   - error: The error that occurred while running the command.
type CommandRunFunc func(p *Program, args []string) error

// Command is a command that can be run by a program.
type Command struct {
	// Name is the name of the command.
	Name string

	// Flags is the list of flags that the command accepts.
	Flags []*Flag[any]

	// Execute is the function that runs the command.
	Execute CommandRunFunc
}

// AddFlag adds a flag to the command.
//
// Parameters:
//   - flag: The flag to add to the command.
func (c *Command) AddFlag(flag *Flag[any]) {
	c.Flags = append(c.Flags, flag)
}
