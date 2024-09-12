package errors

// ErrOrSol is a struct that holds a list of errors and a list of solutions.
type ErrOrSol[T any] struct {
	// error_list is a list of errors.
	error_list []error

	// solution_list is a list of solutions.
	solution_list []T

	// level is the level of the error or solution.
	level int

	// ignore_err is a flag that indicates if the error should be ignored.
	ignore_err bool
}

// AddErr adds an error to the list of errors if the level is greater or equal
// to the current level. Nil errors are ignored.
//
// Parameters:
//   - err: The error to add.
//   - level: The level of the error.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
//
// Behaviors:
//   - If an error has been added with a level greater than the current level,
//     the error list is reset and the new level is updated.
//   - If the error is nil, the ignoreErr flag is set to true and the error list is reset.
func (e *ErrOrSol[T]) AddErr(err error, level int) bool {
	if e == nil {
		return false
	}

	if level < e.level || e.ignore_err || err == nil {
		// Do nothing.
		return true
	}

	if level == e.level {
		e.error_list = append(e.error_list, err)

		return true
	}

	// Clean the previous error list.
	for i := 0; i < len(e.error_list); i++ {
		e.error_list[i] = nil
	}
	e.error_list = nil

	e.error_list = []error{err}
	e.level = level

	return true
}

// AddSol adds a solution to the list of solutions if the level is greater or equal
// to the current level.
//
// Parameters:
//   - sol: The solution to add.
//   - level: The level of the solution.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
//
// Behaviors:
//   - If a solution has been added with a level greater than the current level,
//     the solution list is reset and the new level is updated.
//   - This function sets the ignoreErr flag to true and resets the error list.
func (e *ErrOrSol[T]) AddSol(sol T, level int) bool {
	if e == nil {
		return false
	}

	if level < e.level {
		// Do nothing.
		return true
	}

	if e.level == level {
		e.solution_list = append(e.solution_list, sol)

		return true
	}

	// Clean the previous solution list.

	for i := 0; i < len(e.solution_list); i++ {
		e.solution_list[i] = *new(T)
	}
	e.solution_list = nil

	e.solution_list = []T{sol}
	e.level = level

	if !e.ignore_err {
		e.ignore_err = true

		// Clean the previous error list.
		for i := 0; i < len(e.error_list); i++ {
			e.error_list[i] = nil
		}
		e.error_list = nil
	}

	return true
}

// AddAny adds an element to the list of errors or solutions if the level is greater or equal
// to the current level.
//
// Parameters:
//   - elem: The element to add.
//   - level: The level of the element.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
//
// Behaviors:
//   - If an error has been added with a level greater than the current level,
//     the error list is reset and the new level is updated.
//   - If a solution has been added with a level greater than the current level,
//     the solution list is reset and the new level is updated.
func (e *ErrOrSol[T]) AddAny(elem any, level int) bool {
	if e == nil {
		return false
	}

	if level < e.level || elem == nil {
		// Do nothing.
		return true
	}

	switch elem := elem.(type) {
	case error:
		if e.ignore_err {
			// Do nothing.
			return true
		}

		if level == e.level {
			e.error_list = append(e.error_list, elem)

			return true
		}

		// Clean the previous error list.
		for i := 0; i < len(e.error_list); i++ {
			e.error_list[i] = nil
		}
		e.error_list = nil

		e.error_list = []error{elem}
		e.level = level
	case T:
		if e.level == level {
			e.solution_list = append(e.solution_list, elem)

			return true
		}

		// Clean the previous solution list.
		for i := 0; i < len(e.solution_list); i++ {
			e.solution_list[i] = *new(T)
		}
		e.solution_list = nil

		e.solution_list = []T{elem}
		e.level = level

		if !e.ignore_err {
			e.ignore_err = true

			// Clean the previous error list.
			for i := 0; i < len(e.error_list); i++ {
				e.error_list[i] = nil
			}
			e.error_list = nil
		}
	}

	return true
}

// HasError checks if errors are not ignored and if the error list is not empty.
//
// Returns:
//   - bool: True if errors are not ignored and the error list is not empty, otherwise false.
func (e ErrOrSol[T]) HasError() bool {
	return !e.ignore_err && len(e.error_list) > 0
}

// Errors returns the list of errors. It is a copy of the error list.
//
// Returns:
//   - []error: The list of errors.
func (e ErrOrSol[T]) Errors() []error {
	err_list := make([]error, len(e.error_list))
	copy(err_list, e.error_list)

	return err_list
}

// Solutions returns the list of solutions. It is a copy of the solution list.
//
// Returns:
//   - []T: The list of solutions.
func (e ErrOrSol[T]) Solutions() []T {
	sol_list := make([]T, len(e.solution_list))
	copy(sol_list, e.solution_list)

	return sol_list
}

// Reset resets the ErrOrSol struct to allow for reuse.
func (e *ErrOrSol[T]) Reset() {
	if e == nil {
		return
	}

	if e.error_list != nil {
		for i := 0; i < len(e.error_list); i++ {
			e.error_list[i] = nil
		}
		e.error_list = nil
	}

	if e.solution_list != nil {
		for i := 0; i < len(e.solution_list); i++ {
			e.solution_list[i] = *new(T)
		}
		e.solution_list = nil
	}

	e.level = 0
	e.ignore_err = false
}
