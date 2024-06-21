package common

import (
	"fmt"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

///////////////////////////////////////////////////////

// DescBuilder is a builder for a description.
type DescBuilder struct {
	// lines is the list of lines in the description.
	lines []string
}

// FString implements the ffs.FStringer interface.
func (d *DescBuilder) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	err := trav.AddLines(d.lines)
	if err != nil {
		return err
	}

	return nil
}

// NewDescription creates a new description.
//
// Parameters:
//   - lines: The lines of the description.
//
// Returns:
//   - *Description: The new description.
func NewDescription(lines ...string) *DescBuilder {
	return &DescBuilder{
		lines: lines,
	}
}

// AddLine adds a line to the description.
//
// Parameters:
//   - sections: The sections of the line.
//
// Returns:
//   - *Description: The description.
func (d *DescBuilder) AddLine(sections ...string) *DescBuilder {
	d.lines = append(d.lines, strings.Join(sections, " "))

	return d
}

// Build builds the description.
//
// Returns:
//   - []string: The description.
func (d *DescBuilder) Build() []string {
	linesCopy := make([]string, len(d.lines))
	copy(linesCopy, d.lines)

	d.Reset()

	return linesCopy
}

// Reset resets the description.
func (d *DescBuilder) Reset() {
	for i := range d.lines {
		d.lines[i] = ""
	}

	d.lines = d.lines[:0]
}

type Printer struct {
	lines []string
}

func (p *Printer) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	for _, line := range p.lines {
		err := trav.AddLine(line)
		if err != nil {
			return fmt.Errorf("error adding line: %w", err)
		}
	}

	return nil
}

func NewPrinter(lines []string) *Printer {
	return &Printer{
		lines: lines,
	}
}

type TableAligner struct {
	head    string
	table   [][]string
	tabSize int
}

func NewTableAligner(tabSize int) *TableAligner {
	return &TableAligner{
		table:   make([][]string, 0),
		tabSize: tabSize,
		head:    "",
	}
}

func (ta *TableAligner) SetHead(head string) {
	ta.head = head
}

func (ta *TableAligner) AddRow(row []string) {
	ta.table = append(ta.table, row)
}

func (ta *TableAligner) Reset() {
	ta.table = make([][]string, 0)
	ta.head = ""
}

var (
	tableFunc ffs.FStringFunc[[][]string]
)

func init() {
	tableFunc = func(trav *ffs.Traversor, table [][]string) error {
		if trav == nil {
			return nil
		}

		for _, row := range table {
			trav.AddJoinedLine("", row...)
		}

		return nil
	}
}

func (ta *TableAligner) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	table, err := fs.TabAlign(ta.table, 0, ta.tabSize)
	if err != nil {
		return ue.NewErrWhile("tab aligning", err)
	}

	err = trav.AddLine(ta.head)
	if err != nil {
		return fmt.Errorf("error adding head: %w", err)
	}

	err = ffs.ApplyFormFunc(
		trav.GetConfig(
			ffs.WithModifiedIndent(1),
		),
		trav,
		table,
		tableFunc,
	)
	if err != nil {
		return fmt.Errorf("error applying form: %w", err)
	}

	return nil
}
