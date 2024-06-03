package LyneCml

// Flag is a flag that a command accepts.
type Flag[T any] struct {
	// Flag is the flag that the command accepts.
	Flag string

	// ParseFunc is the function that parses the flag's value.
	ParseFunc func(string) (T, error)

	// Description is the description of the flag.
	Description *Description

	// DefaultValue is the default value of the flag.
	DefaultValue T

	// Value is the value of the flag.
	Value T
}
