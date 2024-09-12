package errors

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// NilValue is the error returned when a pointer is nil. Readers are expected to return this error
	// as is and not wrap it as callers are expected to check for this error using ==.
	NilValue error

	// NilReceiver is the error returned when the receiver is nil. Readers are expected to return this error
	// as is and not wrap it as callers are expected to check for this error using ==.
	NilReceiver error
)

func init() {
	NilValue = errors.New("pointer must not be nil")

	NilReceiver = errors.New("found receiver to be nil")
}

// ErrEmpty represents an error when a value is empty.
type ErrEmpty struct {
	// Type is the type of the empty value.
	Type any
}

// Error implements the error interface.
//
// Message: "{{ .Type }} must not be empty"
func (e ErrEmpty) Error() string {
	var t_string string

	if e.Type == nil {
		t_string = "nil"
	} else {
		to := reflect.TypeOf(e.Type)
		t_string = to.String()
	}

	return t_string + " must not be empty"
}

// NewErrEmpty creates a new ErrEmpty error.
//
// Parameters:
//   - var_type: The type of the empty value.
//
// Returns:
//   - *ErrEmpty: A pointer to the newly created ErrEmpty. Never returns nil.
func NewErrEmpty(var_type any) *ErrEmpty {
	return &ErrEmpty{
		Type: var_type,
	}
}

// ErrGT represents an error when a value is less than or equal to a specified value.
type ErrGT[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be greater than <value>"
func (e ErrGT[T]) Error() string {
	return fmt.Sprintf("value must ge greater than %v", e.Value)
}

// NewErrGT creates a new ErrGT error with the specified value.
//
// Parameters:
//   - value: The minimum value that is not allowed.
//
// Returns:
//   - *ErrGT: A pointer to the newly created ErrGT. Never returns nil.
func NewErrGT[T cmp.Ordered](value T) *ErrGT[T] {
	e := &ErrGT[T]{
		Value: value,
	}
	return e
}

// ErrLT represents an error when a value is greater than or equal to a specified value.
type ErrLT[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be less than <value>"
func (e ErrLT[T]) Error() string {
	return fmt.Sprintf("value must be less than %v", e.Value)
}

// NewErrLT creates a new ErrLT error with the specified value.
//
// Parameters:
//   - value: The maximum value that is not allowed.
//
// Returns:
//   - *ErrLT: A pointer to the newly created ErrLT. Never returns nil.
func NewErrLT[T cmp.Ordered](value T) *ErrLT[T] {
	e := &ErrLT[T]{
		Value: value,
	}
	return e
}

// ErrGTE represents an error when a value is less than a specified value.
type ErrGTE[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be greater than or equal to <value>"
func (e ErrGTE[T]) Error() string {
	return fmt.Sprintf("value must be greater than or equal to %v", e.Value)
}

// NewErrGTE creates a new ErrGTE error with the specified value.
//
// Parameters:
//   - value: The minimum value that is allowed.
//
// Returns:
//   - *ErrGTE: A pointer to the newly created ErrGTE. Never returns nil.
func NewErrGTE[T cmp.Ordered](value T) *ErrGTE[T] {
	e := &ErrGTE[T]{
		Value: value,
	}
	return e
}

// ErrLTE represents an error when a value is greater than a specified value.
type ErrLTE[T cmp.Ordered] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "value must be less than or equal to <value>"
func (e ErrLTE[T]) Error() string {
	return fmt.Sprintf("value must be less than or equal to %v", e.Value)
}

// NewErrLTE creates a new ErrLTE error with the specified value.
//
// Parameters:
//   - value: The maximum value that is allowed.
//
// Returns:
//   - *ErrLTE: A pointer to the newly created ErrLTE. Never returns nil.
func NewErrLTE[T cmp.Ordered](value T) *ErrLTE[T] {
	e := &ErrLTE[T]{
		Value: value,
	}
	return e
}

// ErrUnexpected represents an error when an unexpected value is encountered.
type ErrUnexpected[T any] struct {
	// Expecteds is the list of expected values.
	Expecteds []T

	// Got is the unexpected value.
	Got any

	// Quoted is a flag indicating whether values should be quoted.
	Quoted bool

	// Kind is the kind of the unexpected value.
	Kind string
}

// Error implements the error interface.
//
// Message: "expected <expecteds> <kind>, got <got> instead"
func (e ErrUnexpected[T]) Error() string {
	got := Got(e.Quoted, e.Got)
	elems := StringOfSlice(e.Quoted, e.Expecteds)

	var builder strings.Builder

	builder.WriteString("expected ")

	if len(elems) == 0 {
		if e.Kind != "" {
			builder.WriteString("no ")
			builder.WriteString(e.Kind)
			builder.WriteString(", ")
		} else {
			builder.WriteString("nothing, ")
		}

		builder.WriteString(got)

		return builder.String()
	}

	if len(elems) == 1 {
		builder.WriteString(elems[0])
	} else {
		builder.WriteString("either ")

		if len(elems) > 2 {
			builder.WriteString(strings.Join(elems[:len(elems)-1], ", "))
			builder.WriteRune(',')
		} else {
			builder.WriteString(elems[0])
		}

		builder.WriteString(" or ")
		builder.WriteString(elems[len(elems)-1])
	}

	if e.Kind != "" {
		builder.WriteString(" ")
		builder.WriteString(e.Kind)
	}

	builder.WriteString(", ")
	builder.WriteString(got)

	return builder.String()
}

// NewErrUnexpected creates a new ErrUnexpected error with the specified values.
//
// Parameters:
//   - quoted: A flag indicating whether values should be quoted.
//   - expecteds: The list of expected values.
//   - kind: The kind of the unexpected value.
//   - got: The unexpected value.
//
// Returns:
//   - *ErrUnexpected: A pointer to the newly created ErrUnexpected. Never returns nil.
func NewErrUnexpected[T any](quoted bool, expecteds []T, kind string, got any) *ErrUnexpected[T] {
	return &ErrUnexpected[T]{
		Expecteds: expecteds,
		Got:       got,
		Quoted:    quoted,
		Kind:      kind,
	}
}

// ErrUnexpectedType represents an error when a value has an invalid type.
type ErrUnexpectedType[T any] struct {
	// Elem is the element that caused the error.
	Elem T

	// Kind is the category of the type that was expected.
	Kind string
}

// Error implements the error interface.
//
// Message: "type <type> is not a valid <kind> type"
func (e ErrUnexpectedType[T]) Error() string {
	return fmt.Sprintf("type (%T) is not a valid %s type", e.Elem, e.Kind)
}

// NewErrUnexpectedType creates a new ErrUnexpectedType error.
//
// Parameters:
//   - kind: The name of the type that was expected.
//   - elem: The element that caused the error.
//
// Returns:
//   - *ErrUnexpectedType: A pointer to the newly created ErrUnexpectedType. Never returns nil.
func NewErrUnexpectedType[T any](kind string, elem T) *ErrUnexpectedType[T] {
	e := &ErrUnexpectedType[T]{
		Elem: elem,
		Kind: kind,
	}
	return e
}

// ErrUnexpectedValue represents an error when an unexpected value is encountered.
// This is mostly used in the 'default' case of switch statements.
type ErrUnexpectedValue[T fmt.Stringer] struct {
	// Elem is the element that caused the error.
	Elem T

	// Kind is the category of the type that was expected.
	Kind string
}

// Error implements the error interface.
//
// Message: "<type> (<elem>) is not a supported <kind> type"
func (e ErrUnexpectedValue[T]) Error() string {
	return fmt.Sprintf("%T (%s) is not a supported %s", e.Elem, strconv.Quote(e.Elem.String()), e.Kind)
}

// NewErrUnexpectedValue creates a new ErrUnexpectedValue error.
//
// Parameters:
//   - elem: The element that caused the error.
//   - kind: The name of the type that was expected.
//
// Returns:
//   - *ErrUnexpectedValue: A pointer to the newly created ErrUnexpectedValue. Never returns nil.
func NewErrUnexpectedValue[T fmt.Stringer](elem T, kind string) *ErrUnexpectedValue[T] {
	return &ErrUnexpectedValue[T]{
		Elem: elem,
		Kind: kind,
	}
}

// ErrPanic represents an error when a panic occurs.
type ErrPanic struct {
	// Value is the value that caused the panic.
	Value any
}

// Error implements the error interface.
//
// Message: "panic: {value}"
func (e ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//   - value: The value that caused the panic.
//
// Returns:
//   - *ErrPanic: A pointer to the newly created ErrPanic. Never returns nil.
func NewErrPanic(value any) *ErrPanic {
	return &ErrPanic{
		Value: value,
	}
}

// AssertionFailed is the message that is shown when an assertion fails.
const AssertionFailed string = "[ASSERTION FAILED]: "

// ErrAssertFailed is the error that is shown when an assertion fails.
type ErrAssertFailed struct {
	// Msg is the message that is shown when the assertion fails.
	Msg string
}

// Error implements the error interface.
//
// Message: "[ASSERTION FAILED]: <msg>"
func (e ErrAssertFailed) Error() string {
	var msg string

	if e.Msg == "" {
		msg = "something went wrong"
	} else {
		msg = e.Msg
	}

	return AssertionFailed + msg
}

// NewErrAssertFailed is a constructor for ErrAssertFailed.
//
// Parameters:
//   - msg: the message that is shown when the assertion fails.
//
// Returns:
//   - *ErrAssertFailed: the error. Never returns nil.
func NewErrAssertFailed(msg string) *ErrAssertFailed {
	return &ErrAssertFailed{
		Msg: msg,
	}
}
