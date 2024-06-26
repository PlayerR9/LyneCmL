package Complex

import (
	"errors"
	"fmt"
	"path"
	"strings"

	pd "github.com/PlayerR9/LyneCmL/Complex/display"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
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

	// Options are the optional options of the program.
	Options *Configs

	// display is the display of the program.
	display *pd.Display
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
		return uc.NewErrWhile("tab aligning", err)
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

// Println prints a line to the program's buffer.
//
// Parameters:
//   - args: The items to print.
//
// Returns:
//   - error: An error if the program stopped abruptly.
//     (due to call to Program.Panic())
func (p *Program) Println(args ...interface{}) error {
	str := fmt.Sprintln(args...)
	str = strings.TrimSuffix(str, "\n")

	msg := pd.NewTextMsg(str)

	err := p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Printf prints a formatted line to the program's buffer.
//
// Parameters:
//   - format: The format of the line.
//   - args: The items to print.
//
// Returns:
//   - error: An error if the program stopped abruptly
//     (due to call to Program.Panic())
func (p *Program) Printf(format string, args ...interface{}) error {
	str := fmt.Sprintf(format, args...)

	msg := pd.NewTextMsg(str)

	err := p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// ClearHistory clears the history of the program.
//
// Returns:
//   - error: An error if the history could not be cleared
//     (due to call to Program.Panic())
func (p *Program) ClearHistory() error {
	msg := pd.NewClearHistoryMsg()

	err := p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// SavePartial saves the current history to a file in the partials directory.
//
// This can be used for logging/debugging purposes and/or to save the state of
// the program or evaluate the program's output.
//
// Parameters:
//   - filename: The name of the file to save the partial to.
//
// Returns:
//   - error: An error if the partial could not be saved
//     (due to call to Program.Panic()) or if the filename is empty.
func (p *Program) SavePartial(filename string) error {
	if filename == "" {
		return errors.New("filename is empty")
	}

	fullpath := path.Join("partials", filename)

	msg := pd.NewStoreHistoryMsg(fullpath)

	err := p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Panic causes the program to abruptly exit with the given error.
//
// Parameters:
//   - err: The error that caused the abrupt exit.
//
// Returns:
//   - error: An error if the abrupt exit could not be displayed
//     (due to call to previous Program.Panic() calls).
func (p *Program) Panic(err error) error {
	msg := pd.NewAbruptExitMsg(err)

	err = p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Input requests input from the user.
//
// Parameters:
//   - text: The text to display to the user.
//
// Returns:
//   - string: The input from the user.
//   - error: An error if the input failed.
//
// Errors:
//   - If the context is done.
//   - If input could not be received.
func (p *Program) Input(text string) (string, error) {
	msg := pd.NewInputMsg(text, pd.ItLine)

	err := p.display.Send(msg)
	if err != nil {
		return "", err
	}

	resp, err := msg.Receive()
	if err != nil {
		return "", err
	}

	return resp.(string), nil
}
