package errors

import (
	"fmt"
	"strings"
)

// Enumer is an enum type.
type Enumer interface {
	~int

	// String returns the literal name of the enum.
	//
	// Returns:
	//   - string: The name of the enum.
	String() string
}

// Err is the base error type.
type Err[T Enumer] struct {
	// Code is the error code.
	Code T

	// Message is the error message.
	Message error

	// Suggestions are the error suggestions.
	Suggestions []string
}

// Error implements the error interface.
//
// Message: "{code}: {message}"
func (e Err[T]) Error() string {
	return fmt.Sprintf("%s: %s", e.Code.String(), Error(e.Message))
}

// Unwrap implements errors.Unwrap interface.
func (e Err[T]) Unwrap() error {
	return e.Message
}

// ChangeReason implements the Unwrapper interface.
func (e *Err[T]) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Message = reason

	return true
}

// NewErr creates a new error.
//
// Parameters:
//   - code: The error code.
//   - message: The error message.
//
// Returns:
//   - *Err: The new error. Never returns nil.
func NewErr[T Enumer](code T, message error) *Err[T] {
	return &Err[T]{
		Code:        code,
		Message:     message,
		Suggestions: nil,
	}
}

// AddSuggestion adds a suggestion of the error.
//
// Parameters:
//   - suggestions: The suggestion of the error.
//
// Returns:
//   - *Err: The error. Never returns nil.
//
// Each element in the suggestion is separated by a space but each call to this function
// adds each suggestion on a new line.
func (e *Err[T]) AddSuggestion(suggestions ...string) *Err[T] {
	e.Suggestions = append(e.Suggestions, strings.Join(suggestions, " "))

	return e
}
