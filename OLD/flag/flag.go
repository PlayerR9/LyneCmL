package flag

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ShortFlagPrefix is the prefix for short flags.
	ShortFlagPrefix string = "-"

	// LongFlagPrefix is the prefix for long flags.
	LongFlagPrefix string = "--"
)

// FlagOption is an option for a flag.
//
// Parameters:
//   - flag: The flag to apply the option to.
type FlagOption func(flag *Flag)

// WithShortName sets the short name of the flag.
//
// Parameters:
//   - short_name: The short name of the flag.
//
// Returns:
//   - FlagOption: The option to apply to the flag.
//
// Successive calls to this option will replace the previous short name.
func WithShortName(short_name rune) FlagOption {
	f := func(flag *Flag) {
		flag.short_name = short_name
	}

	return f
}

// WithUsage sets the usage of the flag.
//
// Parameters:
//   - usage: The usage of the flag.
//
// Returns:
//   - FlagOption: The option to apply to the flag.
//
// Successive calls to this option will append the previous usage. Moreover, empty usages will be ignored.
func WithUsage(usage string) FlagOption {
	if usage == "" {
		return func(flag *Flag) {}
	}

	f := func(flag *Flag) {
		flag.usages = append(flag.usages, usage)
	}

	return f
}

// WithDescription sets the description of the flag.
//
// Parameters:
//   - lines: The description of the flag.
//
// Returns:
//   - FlagOption: The option to apply to the flag.
//
// Successive calls to this option will append the previous description. Moreover, empty descriptions will be ignored.
func WithDescription(lines []string) FlagOption {
	if len(lines) == 0 {
		return func(flag *Flag) {}
	}

	f := func(flag *Flag) {
		flag.description = append(flag.description, lines...)
	}

	return f
}

// Valuer is the interface to a value that can be set from a string.
type Valuer interface {
	// Set is a method that sets the value from a string.
	//
	// Parameters:
	//   - str: The string to set the value from.
	//
	// Returns:
	//   - error: An error if the value failed to set.
	Set(str string) error
}

// BoolFlager is the interface that specifies a boolean flag.
type BoolFlager interface {
	// IsBoolFlag is a method that checks if the flag is a boolean flag.
	//
	// Returns:
	//   - bool: True if the flag is a boolean flag. False otherwise.
	IsBoolFlag() bool
}

// Flag is the component of a generic flag.
type Flag struct {
	// long_name is the long name of the flag. (i.e., "--flag")
	long_name string

	// short_name is the short name of the flag. (i.e., "-f")
	short_name rune

	// brief is a brief description of the flag.
	brief string

	// usages are the various usages of the flag.
	usages []string

	// description is a description of the flag.
	description []string

	// value is the current value of the flag after it has been parsed.
	value Valuer
}

// bool_value is a wrapper for bool type.
type bool_value struct {
	// value is the current value of the flag after it has been parsed.
	value bool
}

// Set implements the Valuer interface.
func (b *bool_value) Set(str string) error {
	if str == "" {
		b.value = !b.value

		return nil
	}

	str = strings.ToLower(str)

	switch str {
	case "0", "false", "f":
		b.value = false
	case "1", "true", "t":
		b.value = true
	default:
		return fmt.Errorf("invalid value: %s", str)
	}

	return nil
}

// IsBoolFlag is a method that returns true.
//
// Returns:
//   - bool: True.
func (b *bool_value) IsBoolFlag() bool {
	return true
}

// Bool creates a new bool flag. Panics if the long name is empty.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - def_value: The default value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
//
// Returns:
//   - *bool: A pointer to the boolean value of the flag. Never returns nil.
func Bool(long_name string, def_value bool, brief string, opts ...FlagOption) *bool {
	if long_name == "" {
		panic("long name cannot be empty")
	}

	value := &bool_value{
		value: def_value,
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	flag_set.AddFlag(flag)

	return &value.value
}

// int_value is a wrapper for int type.
type int_value struct {
	// value is the current value of the flag after it has been parsed.
	value int
}

// Set implements the Valuer interface.
func (i *int_value) Set(str string) error {
	num, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	i.value = num

	return nil
}

// Int creates a new int flag. Panics if the long name is empty.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - def_val: The default value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
//
// Returns:
//   - *int: A pointer to the int value of the flag. Never returns nil.
func Int(long_name string, def_val int, brief string, opts ...FlagOption) *int {
	if long_name == "" {
		panic("long name cannot be empty")
	}

	value := &int_value{
		value: def_val,
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	flag_set.AddFlag(flag)

	return &value.value
}

// string_value is a wrapper for string type.
type string_value struct {
	// value is the current value of the flag after it has been parsed.
	value string
}

// Set implements the Valuer interface.
func (s *string_value) Set(str string) error {
	if str == "" {
		return nil
	}

	s.value = str

	return nil
}

// String creates a new string flag. Panics if the long name is empty.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - def_val: The default value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
//
// Returns:
//   - *string: A pointer to the string value of the flag. Never returns nil.
func String(long_name string, def_val string, brief string, opts ...FlagOption) *string {
	if long_name == "" {
		panic("long name cannot be empty")
	}

	value := &string_value{
		value: def_val,
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	flag_set.AddFlag(flag)

	return &value.value
}

// Var creates a new custom flag given a value type. Panics if the long name is empty or the value is nil.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - value: The value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
func Var(long_name string, value Valuer, brief string, opts ...FlagOption) {
	if long_name == "" {
		panic("long name cannot be empty")
	} else if value == nil {
		panic("value cannot be nil")
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	flag_set.AddFlag(flag)
}
