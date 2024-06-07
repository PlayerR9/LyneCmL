package Simple

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

	// Description is a description of the command.
	Description []string

	// Argument is the argument of the command.
	Argument *Argument

	// Run is the function that will be executed when the command is called.
	Run RunFunc
}
