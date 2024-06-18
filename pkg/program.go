package pkg

import (
	"fmt"
	"path"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	llq "github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	ufm "github.com/PlayerR9/MyGoLib/Utility/FileManager"
)

var (
	// ProgramDefaultArgument is the default argument for a program.
	ProgramDefaultArgument *Argument

	// ProgramDefaultRunFunc is the default run function for a program.
	ProgramDefaultRunFunc RunFunc
)

func init() {
	// f is a helper function that parses the arguments of a program.
	//
	// Parameters:
	//   - p: The program to parse the arguments for.
	//   - args: The arguments to parse.
	//
	// Returns:
	//   - *Parsed: The parsed command.
	//   - error: An error if the command could not be parsed.
	f := func(p *Program, args []string) (any, int, error) {
		if len(args) > 0 {
			return nil, 0, fmt.Errorf("unexpected argument: %q", args[0])
		}

		return nil, 0, nil
	}

	ProgramDefaultArgument = &Argument{
		bounds:    [2]int{1, 1},
		parseFunc: f,
	}

	ProgramDefaultRunFunc = func(p *Program, args []string, data any) (any, error) {
		return nil, nil
	}
}

// Program is a program that can be run.
type Program struct {
	*Command

	// commands is a map of commands that the program can execute.
	commands map[string]*Command

	// buffer is a buffer that the program can use to store data.
	buffer *sfb.Buffer[string]

	// history is the history of the program.
	// It is the messages that have been printed to the program.
	history *llq.SafeQueue[string]

	// Version is the version of the program.
	Version string

	// Options are the optional options of the program.
	Options *Configurations
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
	err = trav.AddLine(p.Usage)
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

// ClearHistory clears the history of the program.
//
// This function is thread-safe.
func (p *Program) ClearHistory() {
	p.history.Clear()
}

// SavePartial saves the current history to a file in the partials directory.
//
// This can be used for logging/debugging purposes and/or to save the state of
// the program or evaluate the program's output.
//
// This function is thread-safe.
//
// Parameters:
//   - filename: The name of the file to save the partial to.
//
// Returns:
//   - error: An error if the partial could not be saved.
func (p *Program) SavePartial(filename string) error {
	iter := p.history.Iterator()

	fullpath := path.Join("partials", filename)

	fw := ufm.NewFileWriter(fullpath)

	err := fw.Open()
	if err != nil {
		return err
	}
	defer fw.Close()

	for {
		line, err := iter.Consume()
		if err != nil {
			break
		}

		err = fw.AppendLine(line)
		if err != nil {
			return err
		}
	}

	return nil
}
