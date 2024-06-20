package Simple

import (
	"fmt"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// FlagParseFunc is a function that parses a flag argument.
//
// Parameters:
//   - arg: The argument to parse.
//
// Returns:
//   - T: The parsed argument.
//   - error: An error if the argument failed to parse.
type FlagParseFunc[T any] func(arg string) (T, error)

// Flager is a flag.
type Flager[T any] interface {
	// Fix fixes the flag.
	Fix()

	// GetValue gets the value of the flag.
	//
	// Returns:
	//   - T: The value of the flag.
	GetValue() T

	// SetValue sets the value of the flag.
	//
	// Parameters:
	//   - value: The value to set.
	SetValue(T)

	// GetName gets the name of the flag.
	//
	// Returns:
	//   - string: The name of the flag.
	GetName() string

	CmlComponent
}

// Flag is a flag of type T.
type Flag[T any] struct {
	// Name is the name of the flag.
	Name string

	// Brief is a brief description of the flag.
	Brief string

	// Usage is the usage of the flag.
	Usage string

	// Description is a description of the flag.
	Description *Description

	// DefaultVal is the default value of the flag.
	DefaultVal T

	// ParseFunc is the function that parses the flag argument.
	ParseFunc FlagParseFunc[T]

	// value is the value of the flag.
	value T
}

// GetValue implements Flager interface.
func (f *Flag[T]) GetValue() T {
	return f.value
}

// GetName implements Flager interface.
func (f *Flag[T]) GetName() string {
	return f.Name
}

// SetValue implements Flager interface.
func (f *Flag[T]) SetValue(value T) {
	f.value = value
}

// GenerateUsage implements CmlComponent interface.
//
// Always returns a single string.
func (f *Flag[T]) GenerateUsage() []string {
	var builder strings.Builder

	builder.WriteString(f.Name)
	fmt.Fprintf(&builder, "=<%T>", f.DefaultVal)

	return []string{builder.String()}
}

// Fix implements Flager interface.
func (f *Flag[T]) Fix() {
	f.Name = strings.TrimSpace(f.Name)
	f.Name = "--" + f.Name

	f.Brief = strings.TrimSpace(f.Brief)
	f.Usage = strings.TrimSpace(f.Usage)
}

// FString implements the Flager interface.
func (f *Flag[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	panic("not implemented")
}

// BoolFlag is a special flag that has no value and
// so, it is either present or not.
type BoolFlag struct {
	// Name is the name of the flag.
	Name string

	// Brief is a brief description of the flag.
	Brief string

	// Usage is the usage of the flag.
	Usage string

	// Description is a description of the flag.
	Description *Description

	// value is the value of the flag.
	value bool
}

// GetValue implements Flager interface.
func (bf *BoolFlag) GetValue() bool {
	return bf.value
}

// GetName implements Flager interface.
func (bf *BoolFlag) GetName() string {
	return bf.Name
}

// SetValue implements Flager interface.
func (bf *BoolFlag) SetValue(value bool) {
	bf.value = value
}

// GenerateUsage implements CmlComponent interface.
func (bf *BoolFlag) GenerateUsage() []string {
	return []string{bf.Name}
}

// Fix implements Flager interface.
func (bf *BoolFlag) Fix() {
	bf.Name = strings.TrimSpace(bf.Name)
	bf.Name = "--" + bf.Name
	bf.Brief = strings.TrimSpace(bf.Brief)
	bf.Usage = strings.TrimSpace(bf.Usage)
}

// FString implements the Flager interface.
func (f *BoolFlag) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	panic("not implemented")
}
