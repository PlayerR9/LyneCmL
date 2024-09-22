package simple

import (
	"fmt"
	"strings"
)

var (
	// NoArguments is the default argument.
	NoArguments *Argument
)

func init() {
	NoArguments = &Argument{
		args: nil,
	}
}

// Argument is a struct that represents an argument.
type Argument struct {
	// args is the list of arguments.
	args []string
}

// Fix implements the errors.Fixer interface.
func (a *Argument) Fix() error {
	if a == nil {
		return nil
	}

	return nil
}

// String is a method that returns the string representation of the argument.
//
// Returns:
//   - string: The string representation of the argument.
func (a Argument) String() string {
	elems := make([]string, len(a.args))
	copy(elems, a.args)

	for i := 0; i < len(elems); i++ {
		elems[i] = "<" + elems[i] + ">"
	}

	return strings.Join(elems, " ")
}

// ExactArgs is a helper function that returns an argument with the exact number of arguments.
//
// Parameters:
//   - args: The arguments to return.
//
// Returns:
//   - *Argument: The argument with the exact number of arguments. Never returns nil.
func ExactArgs(args []string) *Argument {
	return &Argument{
		args: args,
	}
}

// parse is a helper function that parses the argument.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - []string: The parsed arguments.
//   - error: An error if the arguments are invalid.
func (a Argument) parse(args []string) ([]string, error) {
	if len(a.args) > len(args) {
		return nil, fmt.Errorf("expected %d arguments, got %d instead", len(a.args), len(args))
	}

	return args[:len(a.args)], nil
}
