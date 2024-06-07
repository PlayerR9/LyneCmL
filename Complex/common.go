package Complex

import (
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// pedanticErrorEval evaluates the error based on the pedantic flag.
//
// Parameters:
//   - isPedantic: The flag that determines whether the program should be pedantic.
//   - reason: The reason for the error.
//   - err: The error to return if the reason is not ignorable.
//
// Returns:
//   - error: The error to return based on the pedantic flag.
func pedanticErrorEval(isPedantic bool, reason, err error) error {
	if reason == nil {
		return nil
	}

	if !ue.As[*ue.ErrIgnorable](reason) || isPedantic {
		return err
	}

	return nil
}
