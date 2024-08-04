package flag

import "strings"

// ErrFlagNotFound is an error that is returned when a flag is not found.
type ErrFlagNotFound struct {
	// IsShort is true if the flag is a short flag. False otherwise.
	IsShort bool

	// Name is the name of the flag.
	Name string
}

// Error implements the error interface.
//
// Message: "flag {{ .Name }} not found"
func (e *ErrFlagNotFound) Error() string {
	var builder strings.Builder

	builder.WriteString("flag ")

	if e.IsShort {
		builder.WriteRune('-')
	} else {
		builder.WriteString("--")
	}

	builder.WriteString(e.Name)

	builder.WriteString(" not found")

	return builder.String()
}

// NewErrFlagNotFound creates a new ErrFlagNotFound.
//
// Parameters:
//   - isShort: True if the flag is a short flag. False otherwise.
//   - name: The name of the flag.
//
// Returns:
//   - *ErrFlagNotFound: The new error. Never returns nil.
func NewErrFlagNotFound(isShort bool, name string) *ErrFlagNotFound {
	e := &ErrFlagNotFound{
		IsShort: isShort,
		Name:    name,
	}
	return e
}

// ErrFlagMissingArg is an error that is returned when a flag is missing an
// argument.
type ErrFlagMissingArg struct {
	// IsShort is true if the flag is a short flag. False otherwise.
	IsShort bool

	// Flag is the flag that is missing an argument.
	Flag *Flag
}

// Error implements the error interface.
//
// Message: "flag {{ .Flag.Name }} requires an argument"
func (e *ErrFlagMissingArg) Error() string {
	var builder strings.Builder

	builder.WriteString("flag ")

	if e.IsShort {
		builder.WriteRune('-')
		builder.WriteRune(e.Flag.short_name)
	} else {
		builder.WriteString("--")
		builder.WriteString(e.Flag.long_name)
	}

	builder.WriteString(" requires an argument")

	return builder.String()
}

// NewErrFlagMissingArg creates a new ErrFlagMissingArg.
//
// Parameters:
//   - isShort: True if the flag is a short flag. False otherwise.
//   - flag: The flag that is missing an argument.
//
// Returns:
//   - *ErrFlagMissingArg: The new error. Never returns nil.
func NewErrFlagMissingArg(isShort bool, flag *Flag) *ErrFlagMissingArg {
	e := &ErrFlagMissingArg{
		IsShort: isShort,
		Flag:    flag,
	}
	return e
}

// ErrInvalidFlag is an error that is returned when an invalid flag is
// passed.
type ErrInvalidFlag struct {
	// IsShort is true if the flag is a short flag. False otherwise.
	IsShort bool

	// Flag is the flag that is invalid.
	Flag *Flag

	// Reason is the reason for the invalid flag.
	Reason error
}

// Error implements the error interface.
//
// Message: "flag {{ .Flag.Name }} is invalid: {{ .Reason }}"
func (e *ErrInvalidFlag) Error() string {
	var builder strings.Builder

	builder.WriteString("flag ")

	if e.Flag == nil {
		builder.WriteString("[no flag specified]")
	} else if e.IsShort {
		builder.WriteRune('-')
		builder.WriteRune(e.Flag.short_name)
	} else {
		builder.WriteString("--")
		builder.WriteString(e.Flag.long_name)
	}

	builder.WriteString(" is invalid")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap is a method that returns the wrapped error.
//
// Returns:
//   - error: The wrapped error.
func (e *ErrInvalidFlag) Unwrap() error {
	return e.Reason
}

// ChangeReason is a method that changes the wrapped error.
//
// Parameters:
//   - reason: The new wrapped error.
func (e *ErrInvalidFlag) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidFlag creates a new ErrInvalidFlag.
//
// Parameters:
//   - isShort: True if the flag is a short flag. False otherwise.
//   - flag: The flag that is invalid.
//   - reason: The reason for the invalid flag.
//
// Returns:
//   - *ErrInvalidFlag: The new error. Never returns nil.
func NewErrInvalidFlag(isShort bool, flag *Flag, reason error) *ErrInvalidFlag {
	e := &ErrInvalidFlag{
		IsShort: isShort,
		Flag:    flag,
		Reason:  reason,
	}
	return e
}
