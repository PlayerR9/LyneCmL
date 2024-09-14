package util

import "strings"

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
