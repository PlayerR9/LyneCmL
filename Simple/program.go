package Simple

import (
	"fmt"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	us "github.com/PlayerR9/MyGoLib/Units/Slice"
)

type Program struct {
	// Name is the name of the program.
	Name string

	// Description is a description of the program.
	Description []string

	// commands is a map of commands that the program can execute.
	commands map[string]*Command

	// buffer is a buffer that the program can use to store data.
	buffer *sfb.Buffer[string]
}

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

func (p *Program) Println(args ...interface{}) {
	str := fmt.Sprintln(args...)

	p.buffer.Send(str)
}

func (p *Program) Printf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)

	p.buffer.Send(str)
}
