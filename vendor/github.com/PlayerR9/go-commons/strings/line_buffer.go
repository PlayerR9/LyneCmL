package strings

import "strings"

// LineBuffer is a struct that represents a line buffer.
type LineBuffer struct {
	// builder is the line builder.
	builder strings.Builder

	// lines is the line buffer.
	lines []string
}

// String returns the lines in the line buffer as a string joined by newlines.
func (lb LineBuffer) String() string {
	if lb.builder.Len() > 0 {
		lb.lines = append(lb.lines, lb.builder.String())
		lb.builder.Reset()
	}

	return strings.Join(lb.lines, "\n")
}

// AddLine adds a line to the line buffer.
//
// Parameters:
//   - line: The line to add.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
func (lb *LineBuffer) AddLine(line string) bool {
	if lb == nil {
		return false
	}

	if lb.builder.Len() > 0 {
		lb.lines = append(lb.lines, lb.builder.String())
		lb.builder.Reset()
	}

	lb.lines = append(lb.lines, line)

	return true
}

// AddString adds a string to the line buffer.
//
// Parameters:
//   - line: The string to add.
//
// Returns:
//   - bool: True if the receiver is not nil, false otherwise.
func (lb *LineBuffer) AddString(line string) bool {
	if lb == nil {
		return false
	}

	lb.builder.WriteString(line)

	return true
}

// Accept accepts the current line buffer.
func (lb *LineBuffer) Accept() {
	if lb == nil || lb.builder.Len() == 0 {
		return
	}

	lb.lines = append(lb.lines, lb.builder.String())
	lb.builder.Reset()
}

// Reset resets the line buffer.
func (lb *LineBuffer) Reset() {
	if lb == nil {
		return
	}

	lb.lines = lb.lines[:0]
	lb.builder.Reset()
}
