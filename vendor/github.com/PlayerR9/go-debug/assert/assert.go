package assert

import (
	"fmt"
)

// Assert panics iff the condition is false. The panic is not a string
// but an error of type *ErrAssertionFailed.
//
// If 'msg' is empty, 'something went wrong' is used insead.
//
// Parameters:
//   - target: the target of the assertion.
//   - cond: the condition to check.
//
// Example:
//
//	foo := "foo"
//	Assert(foo == "bar", "foo is not bar") // panics: "assertion failed: foo is not bar"
func Assert(cond bool, msg string) {
	if cond {
		return
	}

	panic(NewErrAssertFail(msg))
}

// AssertF same as Assert but with a format string and arguments that are in
// accordance with fmt.Printf.
//
// If 'msg' is empty, 'something went wrong' is used insead.
//
// Parameters:
//   - cond: the condition to check.
//   - format: the format string to show if the condition is not met.
//   - args: the arguments to pass to the format string.
//
// Example:
//
//	foo := "foo"
//	bar := "bar"
//	AssertF(foo == bar, "%s is not %s", foo, bar) // panics: "assertion failed: foo is not bar"
func AssertF(cond bool, format string, args ...any) {
	if cond {
		return
	}

	panic(NewErrAssertFail(fmt.Sprintf(format, args...)))
}

// AssertErr is the same as Assert but for errors. Best used for ensuring that a function
// does not return an unexpected error.
//
// Parameters:
//   - err: the error to check.
//   - format: the format describing the function's signature.
//   - args: the arguments passed to the function.
//
// Example:
//
//	foo := "foo"
//	err := my_function(foo, "bar")
//	AssertErr(err, "my_function(%s, %s)", foo, "bar")
//	// panics: "assertion failed: function my_function(foo, bar) returned the error: <err>"
func AssertErr(err error, f_name string, args ...any) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf(f_name, args...) + " = " + err.Error()

	panic(NewErrAssertFail(msg))
}

// AssertOk is the same as Assert but for booleans. Best used for ensuring that a function that
// are supposed to return the boolean `true` does not return `false`.
//
// Parameters:
//   - cond: the result of the function.
//   - format: the format describing the function's signature.
//   - args: the arguments passed to the function.
//
// Example:
//
//	foo := "foo"
//	ok := my_function(foo, "bar")
//	AssertOk(ok, "my_function(%s, %s)", foo, "bar")
//	// panics: "assertion failed: function my_function(foo, bar) returned false while true was expected"
func AssertOk(cond bool, format string, args ...any) {
	if cond {
		return
	}

	msg := fmt.Sprintf(format, args...) + " = false"

	panic(NewErrAssertFail(msg))
}

// AssertDeref tries to dereference an element and panics if it is nil.
//
// Parameters:
//   - elem: the element to dereference.
//   - param_name: the name of the parameter.
//
// Returns:
//   - T: the dereferenced element.
func AssertDeref[T any](elem *T, is_param bool, name string) T {
	if elem != nil {
		return *elem
	}

	var msg string

	if is_param {
		msg = "parameter (" + name + ")"
	} else {
		msg = "variable (" + name + ")"
	}

	msg += " expected to not be nil"

	panic(NewErrAssertFail(msg))
}

// AssertNotNil panics if the element is nil.
//
// Parameters:
//   - elem: the element to check.
//   - param_name: the name of the parameter.
func AssertNotNil(elem any, is_param bool, name string) {
	if elem != nil {
		return
	}

	var msg string

	if is_param {
		msg = "parameter (" + name + ")"
	} else {
		msg = "variable (" + name + ")"
	}

	msg += " expected to not be nil"

	panic(NewErrAssertFail(msg))
}

// AssertTypeOf panics if the element is not of the expected type.
//
// Parameters:
//   - elem: the element to check.
//   - var_name: the name of the variable.
//   - allow_nil: if the element can be nil.
func AssertTypeOf[T any](elem any, target string, allow_nil bool) {
	if elem == nil {
		if !allow_nil {
			msg := fmt.Sprintf("expected %q to be of type %T, got nil instead", target, *new(T))

			panic(NewErrAssertFail(msg))
		}

		return
	}

	_, ok := elem.(T)
	if !ok {
		msg := fmt.Sprintf("expected %q to be of type %T, got %T instead", target, *new(T), elem)

		panic(NewErrAssertFail(msg))
	}
}

// AssertConv tries to convert an element to the expected type and panics if it is not possible.
//
// Parameters:
//   - elem: the element to check.
//   - var_name: the name of the variable.
//
// Returns:
//   - T: the converted element.
func AssertConv[T any](elem any, target string) T {
	if elem == nil {
		msg := fmt.Sprintf("expected %q to be of type %T, got nil instead", target, *new(T))

		panic(NewErrAssertFail(msg))
	}

	res, ok := elem.(T)
	if !ok {
		msg := fmt.Sprintf("expected %q to be of type %T, got %T instead", target, *new(T), elem)

		panic(NewErrAssertFail(msg))
	}

	return res
}