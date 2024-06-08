package OLD

import (
	"strconv"
	"strings"
)

// ErrCommandFailed is an error that occurs when a command fails.
type ErrCommandFailed struct {
	// Command is the command that failed.
	Command string

	// Reason is the reason why the command failed.
	Reason error
}

// Error returns the error message: "command 'command' failed: 'reason'".
//
// Returns:
//   - string: The error message.
func (e *ErrCommandFailed) Error() string {
	var builder strings.Builder

	builder.WriteString("command ")
	builder.WriteString(strconv.Quote(e.Command))
	builder.WriteString(" failed")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// NewErrCommandFailed creates a new ErrCommandFailed.
//
// Parameters:
//   - command: The command that failed.
//   - reason: The reason why the command failed.
//
// Returns:
//   - *ErrCommandFailed: The new ErrCommandFailed.
func NewErrCommandFailed(command string, reason error) *ErrCommandFailed {
	return &ErrCommandFailed{
		Command: command,
		Reason:  reason,
	}
}

// Unwrap returns the reason why the command failed.
//
// Returns:
//   - error: The reason why the command failed.
func (e *ErrCommandFailed) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason why the command failed.
//
// Parameters:
//   - newReason: The new reason why the command failed.
func (e *ErrCommandFailed) ChangeReason(newReason error) {
	e.Reason = newReason
}

// ErrNoCommandProvided is an error that occurs when no command is provided.
type ErrNoCommandProvided struct{}

// Error returns the error message: "no command was provided".
//
// Returns:
//   - string: The error message.
func (e *ErrNoCommandProvided) Error() string {
	return "no command was provided"
}

// NewErrNoCommandProvided creates a new ErrNoCommandProvided.
//
// Returns:
//   - *ErrNoCommandProvided: The new ErrNoCommandProvided.
func NewErrNoCommandProvided() *ErrNoCommandProvided {
	return &ErrNoCommandProvided{}
}

// ErrCommandNotFound is an error that occurs when a command is not found.
type ErrCommandNotFound struct {
	// Command is the command that was not found.
	Command string
}

// Error returns the error message: "command 'command' not found".
//
// Returns:
//   - string: The error message.
func (e *ErrCommandNotFound) Error() string {
	var builder strings.Builder

	builder.WriteString("command ")
	builder.WriteString(strconv.Quote(e.Command))
	builder.WriteString(" not found")

	return builder.String()
}

// NewErrCommandNotFound creates a new ErrCommandNotFound.
//
// Parameters:
//   - command: The command that was not found.
//
// Returns:
//   - *ErrCommandNotFound: The new ErrCommandNotFound.
func NewErrCommandNotFound(command string) *ErrCommandNotFound {
	return &ErrCommandNotFound{
		Command: command,
	}
}

// ErrFailedInitialization is an error that occurs when a program fails to initialize.
type ErrFailedInitialization struct {
	// Reason is the reason why the program failed to initialize.
	Reason error
}

// Error returns the error message: "failed to initialize program: 'reason'".
//
// Returns:
//   - string: The error message.
func (e *ErrFailedInitialization) Error() string {
	var builder strings.Builder

	builder.WriteString("failed to initialize program")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// NewErrFailedInitialization creates a new ErrFailedInitialization.
//
// Parameters:
//   - reason: The reason why the program failed to initialize.
//
// Returns:
//   - *ErrFailedInitialization: The new ErrFailedInitialization.
func NewErrFailedInitialization(reason error) *ErrFailedInitialization {
	return &ErrFailedInitialization{
		Reason: reason,
	}
}

// Unwrap returns the reason why the program failed to initialize.
//
// Returns:
//   - error: The reason why the program failed to initialize.
func (e *ErrFailedInitialization) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason why the program failed to initialize.
//
// Parameters:
//   - newReason: The new reason why the program failed to initialize.
func (e *ErrFailedInitialization) ChangeReason(newReason error) {
	e.Reason = newReason
}
