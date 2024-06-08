package pkg

import "fmt"

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
	return fmt.Sprintf("expected %d arguments, got %d", e.Expected, e.Got)
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
	return fmt.Sprintf("expected at most %d arguments, got %d", e.Expected, e.Got)
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
