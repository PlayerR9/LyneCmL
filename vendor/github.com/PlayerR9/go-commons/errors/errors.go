package errors

import (
	"reflect"
	"strconv"
	"strings"
)

// ErrInvalidParameter represents an error when a parameter is invalid.
type ErrInvalidParameter struct {
	// Parameter is the invalid parameter.
	Parameter string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the error interface.
//
// Message:
// - "parameter (<parameter>) is invalid" if Reason is nil
// - "parameter (<parameter>) is invalid: <reason>" if Reason is not nil
func (e ErrInvalidParameter) Error() string {
	var parameter string

	if e.Parameter != "" {
		parameter = "(" + strconv.Quote(e.Parameter) + ")"
	}

	var builder strings.Builder

	builder.WriteString("parameter ")
	builder.WriteString(parameter)
	builder.WriteString(" is invalid")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the errors.Unwrap interface.
func (e ErrInvalidParameter) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrInvalidParameter) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrInvalidParameter creates a new ErrInvalidParameter error.
//
// Parameters:
//   - parameter: The invalid parameter.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrInvalidParameter: A pointer to the newly created ErrInvalidParameter. Never returns nil.
func NewErrInvalidParameter(parameter string, reason error) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		Parameter: parameter,
		Reason:    reason,
	}
}

// NewErrNilParameter is a convenience method that creates a new *ErrInvalidParameter error
// with a NilValue as the reason.
//
// Parameters:
//   - parameter: The invalid parameter.
//
// Returns:
//   - *ErrInvalidParameter: A pointer to the newly created ErrInvalidParameter. Never returns nil.
func NewErrNilParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		Parameter: parameter,
		Reason:    NilValue,
	}
}

// ErrInvalidUsage represents an error that occurs when a function is used incorrectly.
type ErrInvalidUsage struct {
	// Reason is the reason for the invalid usage.
	Reason error

	// Usage is the usage of the function.
	Usage string
}

// Error is a method of the Unwrapper interface.
//
// Message: "{reason}. {usage}".
//
// However, if the reason is nil, the message is "invalid usage. {usage}" instead.
//
// If the usage is empty, no usage is added to the message.
func (e ErrInvalidUsage) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("invalid usage")
	} else {
		builder.WriteString(e.Reason.Error())
	}

	if e.Usage != "" {
		builder.WriteString(". ")
		builder.WriteString(e.Usage)
	}

	return builder.String()
}

// Unwrap implements the errors.Unwrap interface.
func (e ErrInvalidUsage) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrInvalidUsage) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrInvalidUsage creates a new ErrInvalidUsage error.
//
// Parameters:
//   - reason: The reason for the invalid usage.
//   - usage: The usage of the function.
//
// Returns:
//   - *ErrInvalidUsage: A pointer to the new ErrInvalidUsage error.
func NewErrInvalidUsage(reason error, usage string) *ErrInvalidUsage {
	return &ErrInvalidUsage{
		Reason: reason,
		Usage:  usage,
	}
}

// ErrWhile represents an error that occurs while performing an operation.
type ErrWhile struct {
	// Operation is the operation that was being performed.
	Operation string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "error while {operation}: {reason}"
//
// However, if the reason is nil, the message is "an error occurred while {operation}"
// instead.
func (e ErrWhile) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("an error occurred while ")
		builder.WriteString(e.Operation)
	} else {
		builder.WriteString("error while ")
		builder.WriteString(e.Operation)
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e ErrWhile) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrWhile) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrWhile creates a new ErrWhile error.
//
// Parameters:
//   - operation: The operation that was being performed.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrWhile: A pointer to the newly created ErrWhile.
func NewErrWhile(operation string, reason error) *ErrWhile {
	e := &ErrWhile{
		Operation: operation,
		Reason:    reason,
	}

	return e
}

// ErrAfter is an error that is returned when something goes wrong after a certain
// element in a stream of data.
type ErrAfter struct {
	// After is the element that was processed before the error occurred.
	After string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "after {after}: {reason}".
//
// However, if the reason is nil, the message is "something went wrong after {after}"
// instead.
func (e ErrAfter) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("something went wrong after ")
		builder.WriteString(e.After)
	} else {
		builder.WriteString("after ")
		builder.WriteString(e.After)
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e ErrAfter) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrAfter) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrAfter creates a new ErrAfter error.
//
// Parameters:
//   - after: The element that was processed before the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrAfter: A pointer to the new ErrAfter error. Never returns nil.
func NewErrAfter(after string, reason error) *ErrAfter {
	return &ErrAfter{
		After:  after,
		Reason: reason,
	}
}

// ErrBefore is an error that is returned when something goes wrong before
// a certain element in a stream of data.
type ErrBefore struct {
	// Before is the element that was processed before the error occurred.
	Before string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "before {before}: {reason}".
//
// However, if the reason is nil, the message is "something went wrong before {before}"
// instead.
func (e ErrBefore) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("something went wrong before ")
		builder.WriteString(e.Before)
	} else {
		builder.WriteString("before ")
		builder.WriteString(e.Before)
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrBefore) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrBefore) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrBefore creates a new ErrBefore error.
//
// Parameters:
//   - before: The element that was processed before the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrBefore: A pointer to the new ErrBefore error. Never returns nil.
func NewErrBefore(before string, reason error) *ErrBefore {
	return &ErrBefore{
		Before: before,
		Reason: reason,
	}
}

// ErrInvalidCall represents an error that occurs when a function
// is not called correctly.
type ErrInvalidCall struct {
	// FnName is the name of the function.
	FnName string

	// Signature is the Signature of the function.
	Signature reflect.Type

	// Reason is the Reason for the failure.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "call to {function}({signature}) failed: {reason}".
//
// However, if the reason is nil, the message is "call to {function}({signature})
// failed" instead.
func (e ErrInvalidCall) Error() string {
	var builder strings.Builder

	builder.WriteString("call to ")
	builder.WriteString(e.FnName)
	builder.WriteString(e.Signature.String())
	builder.WriteString(" failed")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrInvalidCall) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrInvalidCall) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrInvalidCall creates a new ErrInvalidCall.
//
// Parameters:
//   - functionName: The name of the function.
//   - function: The function that failed.
//   - reason: The reason for the failure.
//
// Returns:
//   - *ErrInvalidCall: A pointer to the new ErrInvalidCall. Never returns nil.
func NewErrInvalidCall(functionName string, function any, reason error) *ErrInvalidCall {
	return &ErrInvalidCall{
		FnName:    functionName,
		Signature: reflect.ValueOf(function).Type(),
		Reason:    reason,
	}
}
