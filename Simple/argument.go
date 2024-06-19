package Simple

import (
	util "github.com/PlayerR9/LyneCmL/Simple/util"
)

var (
	// NoParseFunc is a function that does nothing.
	NoParseFunc ArgumentParseFunc

	// NoArgument is an argument that takes no arguments.
	NoArgument *Argument
)

func init() {
	NoParseFunc = func(p *Program, args []string) (any, int, error) {
		return nil, len(args), nil
	}

	NoArgument = &Argument{
		bounds:    [2]int{0, 0},
		parseFunc: NoParseFunc,
	}
}

// ArgumentParseFunc is a function that will be executed when the argument is parsed.
//
// Parameters:
//   - p: The program that the argument is being executed on.
//   - args: The arguments that were passed to the argument.
//
// Returns:
//   - any: The result of the arguments.
//   - int: The number of arguments that were parsed.
//   - error: An error if the argument failed to execute.
type ArgumentParseFunc func(p *Program, args []string) (any, int, error)

// Argument is an argument that a command can take.
type Argument struct {
	// bounds is the bounds of the argument.
	bounds [2]int

	// parseFunc is the function that will be executed when the argument is parsed.
	parseFunc ArgumentParseFunc
}

func (a *Argument) fix() {
	if a.parseFunc != nil {
		return
	}

	min, max := a.bounds[0], a.bounds[1]

	if min == 0 {
		if max == 0 {
			// NoArgument
			a.parseFunc = NoParseFunc
		} else {
			// AtMostNArgs
			a.parseFunc = func(p *Program, args []string) (any, int, error) {
				return nil, max, nil
			}
		}
	} else if max == -1 {
		// AtLeastNArgs
		a.parseFunc = func(p *Program, args []string) (any, int, error) {
			return nil, len(args), nil
		}
	} else if min == max {
		// ExactlyNArgs
		a.parseFunc = func(p *Program, args []string) (any, int, error) {
			return nil, min, nil
		}
	} else {
		// RangeArgs
		a.parseFunc = func(p *Program, args []string) (any, int, error) {
			return nil, max, nil
		}
	}
}

// AtLeastNArgs creates an argument that requires at least n arguments.
//
// Parameters:
//   - n: The minimum number of arguments.
//
// Returns:
//   - *Argument: The new argument.
//
// Behaviors:
//   - If n is less than 0, it will be set to 0.
func AtLeastNArgs(n int) *Argument {
	if n < 0 {
		n = 0
	}

	return &Argument{
		bounds:    [2]int{n, -1},
		parseFunc: NoParseFunc,
	}
}

// AtMostNArgs creates an argument that requires at most n arguments.
//
// Parameters:
//   - n: The maximum number of arguments.
//
// Returns:
//   - *Argument: The new argument.
//
// Behaviors:
//   - If n is less than or equal to 0, it will be set to no arguments.
func AtMostNArgs(n int) *Argument {
	if n <= 0 {
		return NoArgument
	} else {
		return &Argument{
			bounds:    [2]int{0, n},
			parseFunc: NoParseFunc,
		}
	}
}

// ExactlyNArgs creates an argument that requires exactly n arguments.
//
// Parameters:
//   - n: The number of arguments.
//
// Returns:
//   - *Argument: The new argument.
//
// Behaviors:
//   - If n is less than or equal to 0, it will be set to no arguments.
func ExactlyNArgs(n int) *Argument {
	if n <= 0 {
		return NoArgument
	} else {
		return &Argument{
			bounds:    [2]int{n, n},
			parseFunc: NoParseFunc,
		}
	}
}

// RangeArgs creates an argument that requires between min and max arguments.
//
// Parameters:
//   - min: The minimum number of arguments.
//   - max: The maximum number of arguments.
//
// Returns:
//   - *Argument: The new argument.
//
// Behaviors:
//   - If min is less than 0, it will be set to 0.
//   - If max is less than 0, it will be set to 0.
//   - If min is greater than max, they will be swapped.
//   - If min and max are both 0, it will be set to no arguments.
func RangeArgs(min, max int) *Argument {
	if min < 0 {
		min = 0
	}

	if max < 0 {
		max = 0
	}

	if min > max {
		min, max = max, min
	}

	if min == 0 && max == 0 {
		return NoArgument
	} else {
		return &Argument{
			bounds:    [2]int{min, max},
			parseFunc: NoParseFunc,
		}
	}
}

// validate is a helper function that validates the number of arguments.
//
// Parameters:
//   - args: The arguments to validate.
//
// Returns:
//   - []string: The arguments if they are valid.
//   - error: An error if the arguments are invalid.
func (a *Argument) validate(args []string) ([]string, error) {
	left := a.bounds[0]

	var right int

	if a.bounds[1] == -1 {
		right = len(args)
	} else {
		right = a.bounds[1]
	}

	if len(args) < left {
		return nil, util.NewErrFewArguments(left, len(args))
	} else if len(args) > right {
		return nil, util.NewErrManyArguments(right, len(args))
	}

	return args, nil
}

// SetParseFunc sets the parse function of the argument.
//
// Parameters:
//   - f: The function to set.
//
// Returns:
//   - *Argument: The argument.
//
// Behaviors:
//   - If f is nil, it will be set to NoParseFunc.
func (a *Argument) SetParseFunc(f ArgumentParseFunc) *Argument {
	if f == nil {
		f = NoParseFunc
	}

	a.parseFunc = f

	return a
}
