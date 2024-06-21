package Simple

import (
	"fmt"
	"strings"

	com "github.com/PlayerR9/LyneCmL/Simple/common"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// Flager is an interface for a flag.
type Flager interface {
	com.Fixer
}

// FlagParseFunc is a function that parses a flag argument.
//
// Parameters:
//   - arg: The argument to parse.
//
// Returns:
//   - any: The parsed argument.
//   - error: An error if the argument failed to parse.
type FlagParseFunc[T any] func(arg string) (T, error)

// Flag is a flag of type any.
type Flag[T any] struct {
	// LongName is the long name of the flag.
	LongName string

	// ShortName is the short name of the flag.
	ShortName rune

	// Brief is a brief description of the flag.
	Brief string

	// Usage is the usage of the flag.
	Usage string

	// Description is a description of the flag.
	Description []string

	// DefaultVal is the default value of the flag.
	DefaultVal T

	// ParseFunc is the function that parses the flag argument.
	ParseFunc FlagParseFunc[T]

	// HasArgument is true if the flag requires an argument. False otherwise.
	HasArgument bool

	// Value is the value of the flag. Do not set this field.
	Value T
}

///////////////////////////////////////////////////////

// Fix implements Flager interface.
func (f *Flag[T]) Fix() {
	f.LongName = strings.TrimSpace(f.LongName)
	f.LongName = "--" + f.LongName

	f.Brief = strings.TrimSpace(f.Brief)
	f.Usage = strings.TrimSpace(f.Usage)
}

// FString implements the FString.FStringer interface.
func (f *Flag[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	panic("not implemented")
}

// GenerateUsage implements CmlComponent interface.
//
// Always returns a single string.
func (f *Flag[T]) GenerateUsage() []string {
	var builder strings.Builder

	builder.WriteString(f.LongName)
	fmt.Fprintf(&builder, "=<%T>", f.DefaultVal)

	return []string{builder.String()}
}
