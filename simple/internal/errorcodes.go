package internal

import (
	"errors"
	"fmt"

	gcers "github.com/PlayerR9/go-commons/errors"
)

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	// InvalidCommand is an error that occurs when a command is invalid.
	InvalidCommand ErrorCode = iota
)

// NewErrMissingCommand is a short-cut for creating an errors.Err error
// with the error code InvalidCommand.
//
// Returns:
//   - *errors.Err[ErrorCode]: The new error. Never returns nil.
func NewErrMissingCommand() *gcers.Err[ErrorCode] {
	err := gcers.NewErr(InvalidCommand, errors.New("no command was provided"))
	err.AddSuggestion("Use 'help' command to display available commands")

	return err
}

// NewErrInvalidCommand is a short-cut for creating an errors.Err error
// with the error code InvalidCommand.
//
// Parameters:
//   - cmd: The command that was invalid.
//
// Returns:
//   - *errors.Err[ErrorCode]: The new error. Never returns nil.
func NewErrInvalidCommand(cmd string) *gcers.Err[ErrorCode] {
	err := gcers.NewErr(InvalidCommand, fmt.Errorf("%q is not a recognized command", cmd))
	err.AddSuggestion("Use 'help' command to display available commands")

	return err
}
