package Simple

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PlayerR9/LyneCml/OLD/Common/util"
)

var (
	// DefaultParseFunc is a function that does nothing.
	DefaultParseFunc ArgumentParseFunc

	// NoArgument is an argument that takes no arguments.
	NoArgument *Argument
)

func init() {
	DefaultParseFunc = func(args []string) (any, error) {
		return args, nil
	}

	NoArgument = &Argument{
		bounds:    [2]int{0, 0},
		parseFunc: DefaultParseFunc,
	}
}

// ArgumentParseFunc is a function that will be executed when the argument is parsed.
//
// Parameters:
//   - args: The arguments that were passed to the argument.
//
// Returns:
//   - any: The result of the arguments.
//   - error: An error if the argument failed to execute.
type ArgumentParseFunc func(args []string) (any, error)

///////////////////////////////////////////////////////

// Argument is an argument that a command can take.
type Argument struct {
	// bounds is the bounds of the argument.
	bounds [2]int

	// parseFunc is the function that will be executed when the argument is parsed.
	parseFunc ArgumentParseFunc
}

// GenerateUsage implements the CmlComponent interface.
func (a *Argument) GenerateUsage() []string {
	min, max := a.bounds[0], a.bounds[1]
	if min == 0 && max == 0 {
		// NoArgument
		return nil
	}

	var lines []string

	if min == 0 {
		// AtMostNArgs
		if max == 1 {
			lines = append(lines, "[arg]")
		} else {
			lines = append(lines, "")

			var builder strings.Builder

			builder.WriteString("(arg1)...(arg")
			builder.WriteString(strconv.Itoa(max))
			builder.WriteRune(')')

			lines = append(lines, builder.String())
		}
	} else if max == -1 {
		// AtLeastNArgs

		if min == 1 {
			lines = append(lines, "(arg) [optional]...")
		} else {
			var builder strings.Builder

			builder.WriteString("(arg1)...(arg")
			builder.WriteString(strconv.Itoa(min))
			builder.WriteRune(')')

			lines = append(lines, builder.String())
		}
	} else if min == max {
		// ExactlyNArgs
		if min <= 2 {
			lines = append(lines, "(arg1) (arg2)")
		} else {
			var builder strings.Builder

			builder.WriteString("(arg1)...(arg")
			builder.WriteString(strconv.Itoa(max))
			builder.WriteRune(')')

			lines = append(lines, builder.String())
		}
	} else {
		// RangeArgs
		var builder strings.Builder

		if min == 1 {
			builder.WriteString("(arg) [optional]...")
		} else {
			builder.WriteString("(arg1)...(arg")
			builder.WriteString(strconv.Itoa(min))
			builder.WriteString(") [optional]...")
		}

		builder.WriteString("...[arg")
		builder.WriteString(strconv.Itoa(max))
		builder.WriteRune(']')

		lines = append(lines, builder.String())
	}

	return lines
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
		parseFunc: DefaultParseFunc,
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
			parseFunc: DefaultParseFunc,
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
			parseFunc: DefaultParseFunc,
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
			parseFunc: DefaultParseFunc,
		}
	}
}

// WithParseFunc sets the parse function of the argument.
//
// Parameters:
//   - f: The function to set.
//
// Returns:
//   - *Argument: The argument.
//
// Behaviors:
//   - If f is nil, it will be set to NoParseFunc.
func (a *Argument) WithParseFunc(f ArgumentParseFunc) *Argument {
	if f == nil {
		f = DefaultParseFunc
	}

	a.parseFunc = f

	return a
}

// ParsedArgument is the result of parsing an argument.
type ParsedArgument struct {
	// CutSet is the cut set of the arguments.
	CutSet []string

	// Data is the data of the arguments.
	Data any

	// Idx is the index of the arguments.
	Idx int
}

// Apply applies the argument to the arguments.
//
// Parameters:
//   - args: The arguments to apply the argument to.
//
// Returns:
//   - *ParsedArgument: The result of the arguments.
//   - error: An error if the argument failed to execute.
func (a *Argument) Apply(args []string) (*ParsedArgument, error) {
	left := a.bounds[0]

	if len(args) < left {
		return nil, util.NewErrFewArguments(left, len(args))
	}

	right := a.bounds[1]

	if right == -1 || right > len(args) {
		right = len(args)
	}

	var data any
	var err error

	for i := right; i >= left; i-- {
		cutSet := args[:i]

		data, err = a.parseFunc(cutSet)
		if err == nil {
			return &ParsedArgument{
				CutSet: cutSet,
				Data:   data,
				Idx:    i,
			}, nil
		}
	}

	return nil, fmt.Errorf("error parsing arguments: %w", err)
}
