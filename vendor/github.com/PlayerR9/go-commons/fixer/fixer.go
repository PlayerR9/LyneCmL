package fixer

import "errors"

var (
	// ValueMustExists is an error indicating that a value must exist. Readers must return this error
	// as it and not wrap it as callers are expected to check for this error using ==.
	ValueMustExists error
)

func init() {
	ValueMustExists = errors.New("value was expected to exist")
}

// Fixer is an interface for objects that can be fixed.
type Fixer interface {
	// Fix fixes the object.
	//
	// Returns:
	//   - error: An error if the object could not be fixed.
	Fix() error
}

// Fix fixes an object if it exists.
//
// Parameters:
//   - elem: The object to fix.
//   - mustExists: A flag indicating if the object must exist.
//
// Returns:
//   - error: An error if the object could not be fixed.
//
// Behaviors:
//   - Returns an error if the object must exist but does not.
//   - Returns nil if the object is nil and does not need to exist.
func Fix(elem Fixer, mustExists bool) error {
	if elem == nil {
		if mustExists {
			return ValueMustExists
		}

		return nil
	}

	err := elem.Fix()
	if err != nil {
		return err
	}

	return nil
}

// Cleaner is an interface for objects that can be cleaned.
type Cleaner interface {
	// Cleanup cleans the object.
	Cleanup()
}

// CleanSlice cleans a slice of elements.
//
// Parameters:
//   - elems: The slice of elements to clean.
//
// Returns:
//   - []T: The cleaned slice of elements.
func CleanSlice[T Cleaner](elems []T) []T {
	if elems == nil {
		return nil
	}

	for i := 0; i < len(elems); i++ {
		current_elem := elems[i]

		current_elem.Cleanup()

		elems[i] = *new(T)
	}

	elems = elems[:0]

	return elems
}

// CleanSliceOf cleans a slice of elements.
//
// Parameters:
//   - elems: The slice of elements to clean.
//
// Returns:
//   - []T: The cleaned slice of elements.
func CleanSliceOf[T any](elems []T) []T {
	if elems == nil {
		return nil
	}

	for i := 0; i < len(elems); i++ {
		elems[i] = *new(T)
	}

	elems = elems[:0]

	return elems
}

// CleanMap cleans a map of elements.
//
// Parameters:
//   - elems: The map of elements to clean.
//
// Returns:
//   - map[K]V: The cleaned map of elements.
func CleanMap[K comparable, V Cleaner](elems map[K]Cleaner) map[K]V {
	if elems == nil {
		return nil
	}

	for k, elem := range elems {
		elem.Cleanup()

		elems[k] = elem

		delete(elems, k)
	}

	return nil
}

// CleanMapOf cleans a map of elements.
//
// Parameters:
//   - elems: The map of elements to clean.
//
// Returns:
//   - map[K]V: The cleaned map of elements.
func CleanMapOf[K comparable, V any](elems map[K]V) map[K]V {
	if elems == nil {
		return nil
	}

	for k := range elems {
		elems[k] = *new(V)

		delete(elems, k)
	}

	return nil
}

// CleanChannel cleans a channel.
//
// Parameters:
//   - ch: The channel to clean.
//
// Returns:
//   - chan T: The cleaned channel.
func CleanChannel[T any](ch chan T) chan T {
	if ch == nil {
		return nil
	}

	close(ch)

	return nil
}

// CleanSimple cleans a simple element.
//
// Parameters:
//   - elem: The element to clean.
//
// Returns:
//   - T: The cleaned element.
func CleanSimple[T comparable](elem T) T {
	return *new(T)
}
