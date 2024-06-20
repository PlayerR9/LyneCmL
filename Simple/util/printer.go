package util

import (
	"fmt"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

type Printer struct {
	lines []string
}

func NewPrinter() *Printer {
	return &Printer{
		lines: make([]string, 0),
	}
}

func (p *Printer) AddLine(line string) {
	p.lines = append(p.lines, line)
}

func (p *Printer) AddJoinedLine(sep string, elems ...string) {
	str := strings.Join(elems, sep)

	p.lines = append(p.lines, str)
}

func (p *Printer) AddEmptyLine() {
	p.lines = append(p.lines, "")
}

func (p *Printer) GetLines() []string {
	return p.lines
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
