package strings

import (
	"slices"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// TableAligner is a struct to help align tables.
type TableAligner struct {
	// head is the head of the table.
	head string

	// table is the table.
	table [][]string

	// idxs is the indexes to align.
	idxs []int
}

// NewTableAligner creates a new table aligner.
func NewTableAligner() *TableAligner {
	return &TableAligner{
		table: make([][]string, 0),
		idxs:  make([]int, 0),
	}
}

// SetHead sets the head of the table.
//
// Parameters:
//   - head: The head of the table.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
func (ta *TableAligner) SetHead(head string) bool {
	if ta == nil {
		return false
	}

	ta.head = head

	return true
}

// AddRow adds a row to the table.
//
// Parameters:
//   - elems: The elements of the row.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
func (ta *TableAligner) AddRow(elems ...string) bool {
	if ta == nil {
		return false
	}

	if len(elems) == 0 {
		ta.table = append(ta.table, []string{""})
	} else {
		ta.table = append(ta.table, elems)
	}

	return true
}

// AlignColumn specifies a row to align.
//
// Parameters:
//   - idx: The index of the row to align.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
//
// If 'idx' is less than zero, this function returns false. If 'idx' already exists
// this function does nothing.
func (ta *TableAligner) AlignColumn(idx int) bool {
	if ta == nil || idx < 0 {
		return false
	}

	pos, ok := slices.BinarySearch(ta.idxs, idx)
	if ok {
		return true
	}

	ta.idxs = slices.Insert(ta.idxs, pos, idx)

	return true
}

// Reset resets the table aligner for reuse.
func (ta *TableAligner) Reset() {
	if ta == nil {
		return
	}

	if len(ta.table) > 0 {
		for i := 0; i < len(ta.table); i++ {
			for j := 0; j < len(ta.table[i]); j++ {
				ta.table[i][j] = ""
			}
			ta.table[i] = ta.table[i][:0]
			ta.table[i] = nil
		}

		ta.table = ta.table[:0]
	}

	ta.head = ""

	if len(ta.idxs) > 0 {
		for i := 0; i < len(ta.idxs); i++ {
			ta.idxs[i] = 0
		}
		ta.idxs = ta.idxs[:0]
	}
}

// Accept accepts the table aligner.
//
// Parameters:
//   - tabSize: The size of the tab.
//   - tableIndent: Whether to indent the table.
//
// Returns:
//   - error: An error of type *errors.ErrInvalidParameter if the
//     tabSize is less than 1.
func (ta TableAligner) Build(tabSize int, tableIndent bool) ([]string, error) {
	if tabSize < 1 {
		return nil, gcers.NewErrInvalidParameter("tabSize", gcers.NewErrGT(0))
	}

	// Add the table indent if needed.
	if tableIndent {
		var builder strings.Builder

		for i := 0; i < len(ta.table); i++ {
			if len(ta.table[i]) == 0 {
				ta.table[i] = append(ta.table[i], "\t")
			} else {
				builder.WriteString("\t")
				builder.WriteString(ta.table[i][0])

				ta.table[i][0] = builder.String()
				builder.Reset()
			}
		}
	}

	// Align the table.
	for _, idx := range ta.idxs {
		ta.table, _ = TabAlign(ta.table, idx, tabSize)
	}

	// Transform the table into a slice of strings.
	var lines []string

	if ta.head != "" {
		lines = append(lines, ta.head)
	}

	for _, row := range ta.table {
		lines = append(lines, strings.Join(row, ""))
	}

	return lines, nil
}
