package table

// fixVerticalBoundaries is a helper function that fixes the vertical boundaries
// of a table.
//
// Parameters:
//   - maxHeight: The maximum height of the table.
//   - elems: The elements of the table.
//   - y: The y-coordinate to fix the boundaries at.
//
// Returns:
//   - [][]T: The elements of the table with the boundaries fixed.
func fixVerticalBoundaries[T any](maxHeight int, elems [][]T, y *int) [][]T {
	actualY := *y

	if actualY < 0 {
		newElems := elems[-actualY:]
		*y = 0

		return newElems
	} else if actualY >= maxHeight {
		newElems := elems[:maxHeight-actualY]
		*y = maxHeight

		return newElems
	}

	totalHeight := len(elems) + actualY
	var newElems [][]T

	if totalHeight <= maxHeight {
		newElems = elems
		y = &totalHeight
	} else {
		newElems = elems[:maxHeight-actualY]
		*y = maxHeight
	}

	return newElems
}

// fixHorizontalBoundaries is a helper function that fixes the horizontal boundaries
// of a table.
//
// Parameters:
//   - maxWidth: The maximum width of the table.
//   - elems: The elements of the table.
//   - x: The x-coordinate to fix the boundaries at.
//
// Returns:
//   - [][]T: The elements of the table with the boundaries fixed.
func fixHorizontalBoundaries[T any](maxWidth int, elems [][]T, x *int) [][]T {
	actualX := *x

	if actualX < 0 {
		for i, row := range elems {
			if -actualX >= len(row) {
				elems[i] = nil
			} else {
				elems[i] = row[-actualX:]
			}
		}

		*x = 0

		return elems
	} else if actualX >= maxWidth {
		for i, row := range elems {
			if actualX >= len(row) {
				elems[i] = nil
			} else {
				elems[i] = row[:maxWidth-actualX]
			}
		}

		*x = maxWidth

		return elems
	}

	for i, row := range elems {
		totalWidth := len(row) + actualX

		if totalWidth < maxWidth {
			continue
		}

		elems[i] = row[:maxWidth-actualX]
	}

	return elems
}

// FixBoundaries is a function that fixes the boundaries of a table of elements based
// on the maximum width and height of the table.
//
// Parameters:
//   - maxWidth: The maximum width of the table.
//   - maxHeight: The maximum height of the table.
//   - elems: The elements of the table.
//   - x: The x-coordinate to fix the boundaries at.
//   - y: The y-coordinate to fix the boundaries at.
//
// Returns:
//   - [][]T: The elements of the table with the boundaries fixed.
//
// Behaviors:
//   - If maxWidth is less than 0, it is set to 0.
//   - If maxHeight is less than 0, it is set to 0.
//   - If elems is empty, nil is returned.
func FixBoundaries[T any](maxWidth, maxHeight int, elems [][]T, x, y *int) [][]T {
	if maxWidth < 0 {
		maxWidth = 0
	}

	if maxHeight < 0 {
		maxHeight = 0
	}

	if len(elems) == 0 {
		*x, *y = 0, 0

		return nil
	}

	elems = fixVerticalBoundaries(maxHeight, elems, y)
	elems = fixHorizontalBoundaries(maxWidth, elems, x)

	return elems
}
