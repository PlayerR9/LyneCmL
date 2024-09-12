package strings

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcint "github.com/PlayerR9/go-commons/ints"
)

// DateStringer prints the date in the format "1st January, 2006".
//
// Parameters:
//
//   - date: The date to print.
//
// Returns:
//
//   - string: The date in the format "1st January, 2006".
func DateStringer(date time.Time) string {
	var builder strings.Builder

	builder.WriteString(gcint.GetOrdinalSuffix(date.Day()))
	builder.WriteRune(' ')
	builder.WriteString(date.Month().String())
	builder.WriteString(", ")
	builder.WriteString(strconv.Itoa(date.Year()))

	return builder.String()
}

// TimeStringer prints the time in the format "3:04 PM".
//
// Parameters:
//
//   - time: The time to print.
//
// Returns:
//
//   - string: The time in the format "3:04 PM".
func TimeStringer(time time.Time) string {
	return time.Format("3:04 PM")
}

// StringsJoiner joins a list of fmt.Stringer values using a separator.
//
// Parameters:
//   - values: The list of fmt.Stringer values to join.
//   - sep: The separator to use when joining the strings.
//
// Returns:
//   - string: The string representation of the values.
func StringsJoiner[T fmt.Stringer](values []T, sep string) string {
	stringValues := make([]string, 0, len(values))

	for _, value := range values {
		stringValues = append(stringValues, value.String())
	}

	return strings.Join(stringValues, sep)
}

// ArrayFormatter formats a list of strings as an array.
//
// Parameters:
//   - values: The list of strings to format.
//
// Returns:
//   - string: The formatted array.
func ArrayFormatter(values []string) string {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteRune(']')

	return builder.String()
}

// findTabStop is a helper function to find the next tab stop for a string.
//
// Parameters:
//   - s: The string to find the tab stop for.
//   - tabSize: The size of the tab.
//
// Returns:
//   - int: The tab stop.
func findTabStop(s string, tabSize int) int {
	s = strings.TrimRight(s, " ")

	count := utf8.RuneCountInString(s)

	return tabSize * ((count / tabSize) + 1)
}

// FixTabStop fixes the tab stops in a string.
//
// The initial level and the integer return are used for
// chaining multiple calls to this function.
//
// Parameters:
//   - init: The initial level of the tab stop.
//   - tabSize: The size of the tab.
//   - spacing: The spacing to use for the tab stop.
//   - str: The string to fix the tab stops for.
//
// Returns:
//   - string: The string with the tab stops fixed.
//   - int: The total number of characters in the string.
//
// Behaviors:
//   - If the tabSize is less than 1, it is set to 1.
//   - If the initial tab stop is less than 0, it is set to 0.
//   - If the string is empty, an empty string is returned.
func FixTabStop(init, tabSize int, spacing, str string) (string, int) {
	var total int

	if init < 0 {
		total = 0
		init = 0
	} else {
		total = init
		init = tabSize - (init % tabSize)
	}

	if len(str) == 0 {
		return "", total
	}

	if tabSize < 1 {
		tabSize = 1
	}

	var builder strings.Builder

	if init > 0 {
		init = tabSize - (init % tabSize)

		for i := init; i > 0; i-- {
			r, size := utf8.DecodeRuneInString(str)
			str = str[size:]

			if r == '\t' {
				repStr := strings.Repeat(spacing, i)
				total += i

				builder.WriteString(repStr)
				break
			} else {
				total++
				builder.WriteRune(r)
			}

			if len(str) == 0 {
				return builder.String(), total
			}
		}
	}

	for {
		for i := tabSize; i > 0; i-- {
			r, size := utf8.DecodeRuneInString(str)
			str = str[size:]

			if r == '\t' {
				repStr := strings.Repeat(spacing, i)
				total += i

				builder.WriteString(repStr)
				break
			} else {
				total++
				builder.WriteRune(r)
			}

			if len(str) == 0 {
				return builder.String(), total
			}
		}
	}
}

// padRight is a helper function to pad a string to the right.
//
// Parameters:
//   - s: The string to pad.
//   - length: The length to pad the string to.
//
// Returns:
//   - string: The padded string.
func padRight(s string, length int) string {
	var builder strings.Builder

	builder.WriteString(s)
	builder.WriteString(strings.Repeat(" ", length-utf8.RuneCountInString(s)))

	return builder.String()
}

// TabAlign aligns the tabs of a table's column.
//
// Parameters:
//   - table: The table to align.
//   - column: The column to align.
//   - tabSize: The size of the tab.
//
// Returns:
//   - [][]string: The aligned table.
//   - error: An error of type *errors.ErrInvalidParameter if the tabSize is less than 1
//     or the column is less than 0.
//
// Behaviors:
//   - If the column is not found in the table, the table is returned as is.
func TabAlign(table [][]string, column int, tabSize int) ([][]string, error) {
	if tabSize < 1 {
		return nil, gcers.NewErrInvalidParameter("tabSize", gcers.NewErrGT(0))
	} else if column < 0 {
		return nil, gcers.NewErrInvalidParameter("column", gcers.NewErrGTE(0))
	}

	seen := make(map[int]bool)

	for i := 0; i < len(table); i++ {
		if len(table[i]) > column {
			seen[i] = true
		}
	}

	if len(seen) == 0 {
		return table, nil
	}

	stops := make(map[int]int)

	for k := range seen {
		table[k][column] = strings.TrimRight(table[k][column], " ")

		stops[k] = findTabStop(table[k][column], tabSize)
	}

	max := -1

	for _, val := range stops {
		if max == -1 || val > max {
			max = val
		}
	}

	for k := range seen {
		table[k][column] = padRight(table[k][column], max)
	}

	return table, nil
}

// TableEntriesAlign aligns the entries of a table.
//
// Parameters:
//   - table: The table to align.
//   - tabSize: The size of the tab.
//
// Returns:
//   - [][]string: The aligned table.
//   - error: An error if there was an issue aligning the table.
//
// Errors:
//   - *errors.ErrAt: If there was an issue aligning a specific column.
//   - *errors.ErrInvalidParameter: If the tabSize is less than 1.
func TableEntriesAlign(table [][]string, tabSize int) ([][]string, error) {
	if tabSize < 1 {
		return nil, gcers.NewErrInvalidParameter("tabSize", gcers.NewErrGT(0))
	}

	width := LongestLine(table)
	if width == -1 {
		return table, nil
	}

	var err error

	for i := 0; i < width; i++ {
		table, err = TabAlign(table, i, tabSize)
		if err != nil {
			return nil, gcint.NewErrAt(i+1, "column", err)
		}
	}

	return table, nil
}

// LongestLine finds the longest line in a table.
//
// Parameters:
//   - table: The table to find the longest line in.
//
// Returns:
//   - int: The length of the longest line. -1 if the table is empty.
func LongestLine[T any](table [][]T) int {
	if len(table) == 0 {
		return -1
	}

	max := -1

	for i := 0; i < len(table); i++ {
		if len(table[i]) > max {
			max = len(table[i])
		}
	}

	return max
}
