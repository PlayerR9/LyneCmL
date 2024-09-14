package simple

import (
	"fmt"
	"strings"
)

const (
	Pipe      string = " | "
	Hellipsis string = "..."
)

var (
	// NoArguments is the default argument.
	NoArguments *Argument
)

func init() {
	NoArguments = &Argument{
		names: []string{},
		min:   0,
		max:   0,
	}
}

// Argument is a struct that represents an argument.
type Argument struct {
	// names is the list of names of the argument.
	names []string

	// min is the minimum number of arguments.
	min int

	// max is the maximum number of arguments. -1 means no maximum.
	max int
}

// String is a method that returns the string representation of the argument.
//
// Returns:
//   - string: The string representation of the argument.
func (a Argument) String() string {
	// [] optional
	// () required
	// | mutually exclusive
	// ... repeating

	if a.min == a.max {
		if a.min == 0 {
			return ""
		}

		elems := make([]string, 0, a.max)

		for i := 0; i < a.max; i++ {
			elems = append(elems, write_arg(a.names[i]))
		}

		return strings.Join(elems, " ")
	}

	if a.max == -1 {
		var builder strings.Builder

		switch a.min {
		case 0:
			builder.WriteRune('[')
			builder.WriteString(write_arg(a.names[0]))
			builder.WriteString(Hellipsis)
			builder.WriteRune(']')
		case 1:
			builder.WriteString(write_arg(a.names[0]))
			builder.WriteString(Hellipsis)
		default:
			builder.WriteString(write_n_args(a.names[0], a.min))
			builder.WriteString(Hellipsis)
		}

		return builder.String()
	}

	if a.min == 0 {
		elems := make([]string, 0, a.max)

		if a.max == 1 {
			elems = append(elems, write_arg(a.names[0]))
		} else {
			for i := 1; i <= a.max; i++ {
				elems = append(elems, write_n_args(a.names[0], i))
			}
		}

		var builder strings.Builder

		builder.WriteRune('[')
		builder.WriteString(strings.Join(elems, Pipe))
		builder.WriteRune(']')

		return builder.String()
	}

	elems := make([]string, 0, a.max-a.min+1)

	for i := a.min; i <= a.max; i++ {
		elems = append(elems, write_n_args(a.names[0], i))
	}

	return strings.Join(elems, Pipe)
}

// AtLeastNArgs is a method that returns an argument that requires at least n arguments.
// If n is less than 0, it will be set to 0.
//
// Parameters:
//   - name: The name of the argument.
//   - n: The minimum number of arguments.
//
// Returns:
//   - *Argument: The argument. Never returns nil.
func AtLeastNArgs(name string, n int) *Argument {
	if n < 0 {
		n = 0
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "arg"
	}

	return &Argument{
		names: []string{name},
		min:   n,
		max:   -1,
	}
}

// AtMostNArgs is a method that returns an argument that requires at most n arguments.
// If n is 0 or less, NoArguments will be returned instead.
//
// Parameters:
//   - name: The name of the argument.
//   - n: The maximum number of arguments.
//
// Returns:
//   - *Argument: The argument. Never returns nil.
func AtMostNArgs(name string, n int) *Argument {
	if n <= 0 {
		return NoArguments
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "arg"
	}

	return &Argument{
		names: []string{name},
		min:   0,
		max:   n,
	}
}

// ExactlyNArgs is a method that returns an argument that requires exactly n arguments.
// If n is 0 or less, NoArguments will be returned instead.
//
// Parameters:
//   - n: The number of arguments.
//
// Returns:
//   - *Argument: The argument. Never returns nil.
func ExactlyNArgs(names []string) *Argument {
	if len(names) == 0 {
		return NoArguments
	}

	for i := 0; i < len(names); i++ {
		names[i] = strings.TrimSpace(names[i])
	}

	var top int

	for i := 0; i < len(names); i++ {
		if names[i] != "" {
			names[top] = names[i]
			top++
		}
	}

	names = names[:top:top]

	if len(names) == 0 {
		return NoArguments
	}

	return &Argument{
		names: names,
		min:   len(names),
		max:   len(names),
	}
}

// BetweenArgs is a method that returns an argument that requires between min and max arguments.
// If min is less than 0, it will be set to 0.
// If max is less than 0, it will be set to 0.
// If min is greater than max, min and max will be swapped.
//
// Parameters:
//   - name: The name of the argument.
//   - min: The minimum number of arguments.
//   - max: The maximum number of arguments.
//
// Returns:
//   - *Argument: The argument. Never returns nil.
func BetweenArgs(name string, min, max int) *Argument {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "arg"
	}

	if min < 0 {
		min = 0
	}

	if max < 0 {
		max = 0
	}

	if min >= max {
		min, max = max, min
	}

	return &Argument{
		names: []string{name},
		min:   min,
		max:   max,
	}
}

// Check is a helper function that checks the argument.
//
// Parameters:
//   - args: The arguments to check.
//
// Returns:
//   - []string: The checked arguments.
//   - error: An error if the arguments are invalid.
func (a Argument) Check(args []string) ([]string, error) {
	if len(args) < a.min {
		return nil, fmt.Errorf("expected at least %d arguments, got %d instead", a.min, len(args))
	}

	var max int

	if a.max == -1 || a.max > len(args) {
		max = len(args)
	} else {
		max = a.max
	}

	return args[:max:max], nil
}

// write_arg is a helper function that writes an argument.
//
// Parameters:
//   - name: The name of the argument.
//
// Returns:
//   - string: The string representation of the argument.
func write_arg(name string) string {
	return "<" + name + ">"
}

// write_n_args is a helper function that writes n arguments.
//
// Parameters:
//   - n: The number of arguments.
//
// Returns:
//   - string: The string representation of the arguments.
//
// Assertions:
//   - n >= 0
func write_n_args(name string, n int) string {
	if n == 0 {
		return ""
	}

	elems := make([]string, 0, n)

	for i := 0; i < n; i++ {
		elems = append(elems, write_arg(name))
	}

	return strings.Join(elems, " ")
}
