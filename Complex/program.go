package Complex

import (
	"context"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	llq "github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// FilterInvalidCmd is a filter that filters out invalid commands.
	FilterInvalidCmd us.PredicateFilter[*Command]
)

func init() {
	FilterInvalidCmd = func(cmd *Command) bool {
		if cmd == nil {
			return false
		}

		cmd.Name = strings.TrimSpace(cmd.Name)
		return cmd.Name != ""
	}
}

// Program is a program that can be run.
type Program struct {
	// Name is the name of the command.
	Name string

	// Brief is a brief description of the command.
	Brief string

	// Description is a description of the command.
	Description *Description

	// Version is the version of the program.
	Version string

	// commands is a map of commands that the program can execute.
	commands map[string]*Command

	// buffer is a buffer that the program can use to store data.
	buffer *sfb.Buffer[any]

	// history is the history of the program.
	// It is the messages that have been printed to the program.
	history *llq.SafeQueue[string]

	// Options are the optional options of the program.
	Options *Configurations

	// ctx is the context of the program.
	ctx context.Context
}

// GenerateProgram creates a new program with the given command, version, options,
// and subcommands.
//
// Parameters:
//   - cmd: The main/privileged command of the program.
//   - version: The version of the program.
//   - opts: The options of the program.
//   - subCmds: The subcommands of the program.
//
// Returns:
//   - *Program: The new program.
//
// Behaviors:
//   - If cmd is nil, it will be set to a new command.
//   - nil subcommands will be filtered out.
//   - If a subcommand has the same name as another subcommand, the first one
//     will be kept.
func (p *Program) fix() {
	p.Name = strings.TrimSpace(p.Name)
	p.Brief = strings.TrimSpace(p.Brief)

	if p.Options == nil {
		p.Options = DefaultOptions
	} else {
		p.Options.fix()
	}

	if p.commands == nil {
		p.commands = make(map[string]*Command)
	}

	p.commands[HelpCmdOpcode] = HelpCmd
}

// SetCommands sets the commands of the program.
//
// Parameters:
//   - cmds: The commands to set.
//
// Behaviors:
//   - nil commands will be filtered out.
//   - If a command has the same name as another command, the first one
//     will be kept.
//   - The Help command will overwrite any other command with the same name.
func (p *Program) SetCommands(cmds ...*Command) {
	cmds = us.SliceFilter(cmds, FilterInvalidCmd)
	if len(cmds) == 0 {
		return
	}

	if p.commands == nil {
		p.commands = make(map[string]*Command)
	}

	for _, cmd := range cmds {
		cmd.fix()

		_, ok := p.commands[cmd.Name]
		if !ok {
			p.commands[cmd.Name] = cmd
		}
	}
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

// GetSpacing gets the spacing between columns.
//
// Returns:
//   - int: The spacing between columns.
func (p *Program) GetSpacing() int {
	return p.Options.Spacing
}

// DisplayHelp displays the help of the program.
//
// Returns:
//   - []string: The lines of the help.
//   - error: An error if the help could not be displayed.
func (p *Program) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	var err error

	// Program: <name> - <brief>
	if p.Brief == "" {
		err = trav.AddJoinedLine(" ", "Program:", p.Name)
	} else {
		err = trav.AddJoinedLine(" ", "Program:", p.Name, "-", p.Brief)
	}
	if err != nil {
		return err
	}

	trav.EmptyLine()

	// Version: <version>
	if p.Version != "" {
		err := trav.AddJoinedLine(" ", "Version:", p.Version)
		if err != nil {
			return err
		}

		trav.EmptyLine()
	}

	// Usage: <name> (command) [arguments]
	err = trav.AddJoinedLine(" ", "Usage:", p.Name, "(command)", "[arguments]")
	if err != nil {
		return err
	}

	trav.EmptyLine()

	// Description:
	// 	<description>
	if p.Description != nil {
		err := trav.AddLine("Description:")
		if err != nil {
			return err
		}

		err = ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithModifiedIndent(1),
			),
			trav,
			p.Description,
		)
		if err != nil {
			return err
		}

		trav.EmptyLine()
	}

	// Commands:
	// 	<usage> 	 <brief> (vertical alignment)
	//    ...
	table := make([][]string, 0, len(p.commands))
	for _, command := range p.commands {
		table = append(table, []string{command.Usage, command.Brief})
	}

	table, err = fs.TabAlign(table, 0, p.GetTabSize())
	if err != nil {
		return ue.NewErrWhile("tab aligning", err)
	}

	err = trav.AddLine("Commands:")
	if err != nil {
		return err
	}

	err = ffs.ApplyFormFunc(
		trav.GetConfig(
			ffs.WithModifiedIndent(1),
		),
		trav,
		table,
		func(trav *ffs.Traversor, table [][]string) error {
			if trav == nil {
				return nil
			}

			for _, row := range table {
				trav.AddJoinedLine("", row...)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
