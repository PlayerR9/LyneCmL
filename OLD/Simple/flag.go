package Simple

import (
	"fmt"
	"strings"
	// ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

const (
	// ShortFlagPrefix is the prefix for short flags.
	ShortFlagPrefix string = "-"

	// LongFlagPrefix is the prefix for long flags.
	LongFlagPrefix string = "--"
)

var (
	// DefaultFlagArgument is the default flag argument.
	// That is, a flag that returns the argument as a string and
	// whose default value is an empty string.
	DefaultFlagArgument *FlagArgument
)

func init() {
	DefaultFlagArgument = &FlagArgument{
		defaultVal: "",
		parseFunc: func(arg string) (any, error) {
			return arg, nil
		},
	}
}

// FlagParseFunc is a function that parses a flag argument.
//
// Parameters:
//   - arg: The argument to parse.
//
// Returns:
//   - any: The parsed argument.
//   - error: An error if the argument failed to parse.
type FlagParseFunc func(arg string) (any, error)

// FlagArgument is a flag argument.
type FlagArgument struct {
	// defaultVal is the default value of the flag.
	defaultVal any

	// parseFunc is the function that parses the flag argument.
	parseFunc FlagParseFunc
}

// Flag is a flag that a program can use.
type Flag struct {
	// LongName is the long name of the flag.
	LongName string

	// ShortName is the short name of the flag.
	ShortName rune

	// Brief is a brief description of the flag.
	Brief string

	// Usages are the usages of the flag.
	Usages []string

	// Description is a description of the flag.
	Description []string

	// value is the value of the flag.
	value any

	// Argument is the Argument of the flag.
	Argument *FlagArgument
}

// NewFlagArgument creates a new flag argument.
//
// Parameters:
//   - defaultVal: The default value of the flag.
//   - parseFunc: The function that parses the flag argument.
//
// Returns:
//   - *FlagArgument: The new flag argument.
//
// Behaviors:
//   - If parseFunc is nil, then the flag is set to the default flag argument.
func NewFlagArgument(defaultVal any, parseFunc FlagParseFunc) *FlagArgument {
	var fa *FlagArgument

	if parseFunc == nil {
		fa = DefaultFlagArgument
	} else {
		fa = &FlagArgument{
			defaultVal: defaultVal,
			parseFunc:  parseFunc,
		}
	}

	return fa
}

// NewNoFlagArgument creates a new flag argument that does not take an argument.
//
// Returns:
//   - *FlagArgument: The new flag argument.
func NewNoFlagArgument() *FlagArgument {
	return nil
}

// Apply applies the flag to the argument.
//
// Parameters:
//   - arg: The argument to apply.
//
// Returns:
//   - error: An error if the flag failed to apply.
func (f *Flag) Apply(arg string) error {
	if f.Argument == nil {
		f.value = true
	} else {
		value, err := f.Argument.parseFunc(arg)
		if err != nil {
			return err
		}

		f.value = value
	}

	return nil
}

// Value gets the value of the flag.
//
// Returns:
//   - any: The value of the flag.
func (f *Flag) Value() any {
	return f.value
}

///////////////////////////////////////////////////////

/* // FString implements the FString.FStringer interface.
func (f *Flag) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	panic("not implemented")
} */

// GenerateUsage implements CmlComponent interface.
//
// Always returns at most two lines.
func (f *Flag) GenerateUsage() []string {
	var lines []string

	var builder strings.Builder

	if f.ShortName != 0 {
		builder.WriteString(ShortFlagPrefix)
		builder.WriteRune(f.ShortName)

		if f.Argument != nil {
			fmt.Fprintf(&builder, "=<%T>", f.Argument.defaultVal)
		}

		lines = append(lines, builder.String())
		builder.Reset()
	}

	builder.WriteString(LongFlagPrefix)
	builder.WriteString(f.LongName)

	if f.Argument != nil {
		fmt.Fprintf(&builder, "=<%T>", f.Argument.defaultVal)
	}

	lines = append(lines, builder.String())

	return lines
}
