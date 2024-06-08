package Simple

import (
	"fmt"
	"strings"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	us "github.com/PlayerR9/MyGoLib/Units/Slice"
)

// Program is a program that can be run.
type Program struct {
	// Name is the name of the program.
	Name string

	// Brief is a brief description of the program.
	Brief string

	// Description is a description of the program.
	Description []string

	// commands is a map of commands that the program can execute.
	commands map[string]*Command

	// buffer is a buffer that the program can use to store data.
	buffer *sfb.Buffer[string]

	// Options are the optional options of the program.
	Options *ProgramOptions
}

// AddCommands adds commands to the program.
//
// Parameters:
//   - commands: The commands to add to the program.
//
// Behaviors:
//   - If commands is empty, no commands will be added.
//   - Nil commands will be filtered out.
//   - Commands with the same name will overwrite existing commands.
//   - The help command overwrites anything so, never specify it.
func (p *Program) AddCommands(commands ...*Command) {
	commands = us.FilterNilValues(commands)
	if len(commands) == 0 {
		return
	}

	if p.commands == nil {
		p.commands = make(map[string]*Command)
	}

	for _, command := range commands {
		p.commands[command.Name] = command
	}
}

// Println prints a line to the program's buffer.
//
// Parameters:
//   - args: The items to print.
func (p *Program) Println(args ...interface{}) {
	str := fmt.Sprintln(args...)

	p.buffer.Send(str)
}

// Printf prints a formatted line to the program's buffer.
//
// Parameters:
//   - format: The format of the line.
//   - args: The items to print.
//
// Behaviors:
//   - A newline character will be appended to the end of the string
//     if it does not already have one.
func (p *Program) Printf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)

	p.buffer.Send(str)
}

// GetTabSize gets the size of a tab character.
//
// Returns:
//   - int: The size of a tab character.
func (p *Program) GetTabSize() int {
	return p.Options.TabSize
}

// GetTab gets a string of tabs.
//
// Returns:
//   - string: The tab string.
func (p *Program) GetTab() string {
	return strings.Repeat(" ", p.Options.TabSize)
}

// fix is a helper function that fixes the program in order to
// make it easier to use.
func (p *Program) fix(arg string) {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		p.Name = arg
	}

	p.Brief = strings.TrimSpace(p.Brief)

	if p.Options == nil {
		p.Options = DefaultOptions
	}

	newCommands := make(map[string]*Command)

	for _, command := range p.commands {
		command.fix()

		if command.Name == "" {
			continue
		}

		newCommands[command.Name] = command
	}

	p.commands = newCommands

	p.Options.fix()
}
