package fixer

import (
	"strconv"
	"strings"
)

// ErrFix is an error indicating that a field could not be fixed.
type ErrFix struct {
	// Field is the field that could not be fixed.
	Field string

	// Reason is the reason the field could not be fixed.
	Reason error
}

// Error implements the errors.Unwrapper interface.
//
// Message:
//   - "failed to fix field <field>" if the reason is nil.
//   - "field <field> failed to fix: <reason>" if the reason is not nil.
func (e ErrFix) Error() string {
	var values []string

	if e.Reason == nil {
		values = []string{
			"failed to fix field",
			strconv.Quote(e.Field),
		}
	} else {
		values = []string{
			"field",
			strconv.Quote(e.Field),
			"failed to fix:",
			e.Reason.Error(),
		}
	}

	msg := strings.Join(values, " ")

	return msg
}

// Unwrap implements the errors.Unwrapper interface.
func (e ErrFix) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the errors.Unwrapper interface.
func (e *ErrFix) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrFix creates a new ErrFix error.
//
// Parameters:
//   - field: The field that could not be fixed.
//   - reason: The reason the field could not be fixed.
//
// Returns:
//   - *ErrFix: The new error. Never returns nil.
func NewErrFix(field string, reason error) *ErrFix {
	e := &ErrFix{
		Field:  field,
		Reason: reason,
	}
	return e
}

// ErrFixAt is an error indicating that a field at an index could not be fixed.
type ErrFixAt struct {
	// Field is the field that could not be fixed.
	Field string

	// Idx is the index of the field that could not be fixed.
	Idx int

	// Reason is the reason the field could not be fixed.
	Reason error
}

// Error implements the errors.Unwrapper interface.
//
// Message:
//   - "failed to fix field <field> at index <idx>" if the reason is nil.
//   - "field <field> at index <idx> failed to fix: <reason>" if the reason is not nil.
func (e ErrFixAt) Error() string {
	var values []string

	if e.Reason == nil {
		values = []string{
			"failed to fix field",
			strconv.Quote(e.Field),
			"at index",
			"(",
			strconv.Itoa(e.Idx),
			")",
		}
	} else {
		values = []string{
			"field",
			strconv.Quote(e.Field),
			"at index",
			"(",
			strconv.Itoa(e.Idx),
			")",
			"failed to fix:",
			e.Reason.Error(),
		}
	}

	msg := strings.Join(values, " ")

	return msg
}

// Unwrap implements the errors.Unwrapper interface.
func (e *ErrFixAt) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the errors.Unwrapper interface.
func (e *ErrFixAt) ChangeReason(reason error) bool {
	if e == nil {
		return false
	}

	e.Reason = reason

	return true
}

// NewErrFixAt creates a new ErrFixAt error.
//
// Parameters:
//   - field: The field that could not be fixed.
//   - idx: The index of the field that could not be fixed.
//   - reason: The reason the field could not be fixed.
//
// Returns:
//   - *ErrFixAt: The new error. Never returns nil.
func NewErrFixAt(field string, idx int, reason error) *ErrFixAt {
	e := &ErrFixAt{
		Field:  field,
		Idx:    idx,
		Reason: reason,
	}

	return e
}
