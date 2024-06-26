package util

import (
	"strconv"
	"strings"
)

// ErrFlagConflict is an error for when a flag conflicts with another flag.
type ErrFlagConflict struct {
	// Flag is the flag that conflicts with another flag.
	Flag string

	// Conflicting is the flag that conflicts with the flag.
	Conflicting string

	// Reason is the reason for the conflict.
	Reason error
}

// Error implements the errors.Unwrapper interface.
//
// Message:
//   - "flag --<Flag> conflicts with flag --<Conflicting>: <Reason>"
//   - "flag --<Flag> conflicts with flag --<Conflicting>" if Reason is nil.
func (e *ErrFlagConflict) Error() string {
	var builder strings.Builder

	builder.WriteString("flag --")
	builder.WriteString(e.Flag)
	builder.WriteString(" conflicts with flag --")
	builder.WriteString(e.Conflicting)

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the errors.Unwrapper interface.
func (e *ErrFlagConflict) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the errors.Unwrapper interface.
func (e *ErrFlagConflict) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrFlagConflict creates a new ErrFlagConflict.
//
// Parameters:
//   - flag: The flag that conflicts with another flag.
//   - conflicting: The flag that conflicts with the flag.
//   - reason: The reason for the conflict.
//
// Returns:
//   - *ErrFlagConflict: The new ErrFlagConflict.
func NewErrFlagConflict(flag, conflicting string, reason error) *ErrFlagConflict {
	e := &ErrFlagConflict{
		Flag:        flag,
		Conflicting: conflicting,
		Reason:      reason,
	}
	return e
}

///////////////////////////////////////////////////////

// ErrFewArguments is an error that is returned when too few arguments are passed.
type ErrFewArguments struct {
	// Expected is the number of arguments that were expected.
	Expected int

	// Got is the number of arguments that were passed.
	Got int
}

// Error implements the error interface.
//
// Returns the message: "expected 'expected' arguments, got 'got'".
func (e *ErrFewArguments) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.Itoa(e.Expected))
	builder.WriteString(" arguments, got ")
	builder.WriteString(strconv.Itoa(e.Got))
	builder.WriteString(" instead")

	return builder.String()
}

// NewErrFewArguments creates a new ErrFewArguments.
//
// Parameters:
//   - expected: The number of arguments that were expected.
//   - got: The number of arguments that were passed.
//
// Returns:
//   - *ErrFewArguments: The new ErrFewArguments.
func NewErrFewArguments(expected, got int) *ErrFewArguments {
	return &ErrFewArguments{
		Expected: expected,
		Got:      got,
	}
}

// ErrManyArguments is an error that is returned when too many arguments are passed.
type ErrManyArguments struct {
	// Expected is the number of arguments that were expected.
	Expected int

	// Got is the number of arguments that were passed.
	Got int
}

// Error implements the error interface.
//
// Returns the message: "expected at most 'expected' arguments, got 'got'".
func (e *ErrManyArguments) Error() string {
	var builder strings.Builder

	builder.WriteString("expected at most ")
	builder.WriteString(strconv.Itoa(e.Expected))
	builder.WriteString(" arguments, got ")
	builder.WriteString(strconv.Itoa(e.Got))
	builder.WriteString(" instead")

	return builder.String()
}

// NewErrManyArguments creates a new ErrManyArguments.
//
// Parameters:
//   - expected: The number of arguments that were expected.
//   - got: The number of arguments that were passed.
//
// Returns:
//   - *ErrManyArguments: The new ErrManyArguments.
func NewErrManyArguments(expected, got int) *ErrManyArguments {
	return &ErrManyArguments{
		Expected: expected,
		Got:      got,
	}
}

// ErrNoCommand is an error that is returned when no command is provided.
type ErrNoCommand struct{}

// Error implements the error interface.
//
// Returns the message: "no command provided".
func (e *ErrNoCommand) Error() string {
	return "no command provided"
}

// NewErrNoCommand creates a new ErrNoCommand.
//
// Returns:
//   - *ErrNoCommand: The new ErrNoCommand.
func NewErrNoCommand() *ErrNoCommand {
	return &ErrNoCommand{}
}

// ErrUnknownCommand is an error that is returned when an unknown command is provided.
type ErrUnknownCommand struct {
	// Command is the unknown command.
	Command string
}

// Error implements the error interface.
//
// Returns the message: "command 'command' not found".
func (e *ErrUnknownCommand) Error() string {
	var builder strings.Builder

	builder.WriteString("sub command ")
	builder.WriteString(strconv.Quote(e.Command))
	builder.WriteString(" not found")

	return builder.String()
}

// NewErrUnknownCommand creates a new ErrUnknownCommand.
//
// Parameters:
//   - command: The unknown command.
//
// Returns:
//   - *ErrUnknownCommand: The new ErrUnknownCommand.
func NewErrUnknownCommand(command string) *ErrUnknownCommand {
	return &ErrUnknownCommand{
		Command: command,
	}
}
