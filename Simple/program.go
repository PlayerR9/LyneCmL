package Simple

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	pd "github.com/PlayerR9/LyneCmL/Simple/display"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	utfm "github.com/PlayerR9/MyGoLib/Utility/FileManager"
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
	Description []string

	// Version is the version of the program.
	Version string

	// Configs is the configurations of the program.
	Configs *pd.DisplayConfigs

	// commands is a map of commands that the program can execute.
	commands map[string]*Command

	// display is the display of the program.
	display *pd.Display

	// flags is a map of flags that the program can use.
	flags map[string]*Flag
}

// GetCmdMap gets the map of commands of the program.
//
// Returns:
//   - map[string]*Command: The map of commands.
func (p *Program) GetCmdMap() map[string]*Command {
	return p.commands
}

// SetDisplay sets the display of the program.
//
// No operation if display is nil.
//
// Parameters:
//   - display: The display to set.
func (p *Program) SetDisplay(display *pd.Display) {
	if display == nil {
		return
	}

	p.display = display
}

// GenerateUsage implements the CmlComponent interface.
//
// Always one line.
func (p *Program) GenerateUsage() []string {
	var builder strings.Builder

	builder.WriteString(p.Name)
	builder.WriteString(" (command) [arguments] [flags]")

	return []string{builder.String()}
}

// LoadConfigs loads the configuration from the file.
//
// Returns:
//   - error: An error if the configuration failed to load.
func (p *Program) LoadConfigs() error {
	ok, err := utfm.FileExists(pd.ConfigLoc)
	if err != nil {
		return err
	}

	f, err := utfm.Create(pd.ConfigLoc, utfm.DP_All, utfm.FP_All)
	if err != nil {
		return err
	}
	defer f.Close()

	if !ok {
		// Write the configuration to the file
		data, err := json.MarshalIndent(p.Configs, "", "  ")
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		if err != nil {
			return err
		}
	} else {
		stat, err := f.Stat()
		if err != nil {
			return err
		}

		data := make([]byte, stat.Size())

		_, err = f.Read(data)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &p.Configs)
		if err != nil {
			return err
		}

		// Fix the configurations
		err = p.Configs.Fix()
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveConfigs saves the configuration to the file.
//
// Returns:
//   - error: An error if the configuration failed to save.
func (p *Program) SaveConfigs() error {
	data, err := json.MarshalIndent(p.Configs, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(pd.ConfigDir, utfm.DP_All)
	if err != nil {
		return err
	}

	err = os.WriteFile(pd.ConfigLoc, data, utfm.FP_All)
	if err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////

// FString implements the CmlComponent interface.
func (p *Program) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	settings := &CommandFSSetting{
		Spacing: " ",
	}

	for _, opt := range opts {
		opt(settings)
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
	usage := p.GenerateUsage()

	err = trav.AddJoinedLine(" ", "Usage:", usage[0])
	if err != nil {
		return err
	}

	trav.EmptyLine()

	ta := fs.NewTableAligner()
	tabSize := trav.GetConfig().GetTabSize()

	// Description:
	// 	<description>
	if p.Description != nil {
		ta.SetHead("Description:")

		for _, desc := range p.Description {
			ta.AddRow(desc)
		}

		lines, _ := ta.Build(tabSize, true)

		err = trav.AddLines(lines)
		if err != nil {
			return err
		}

		trav.EmptyLine()
	}

	// Commands:
	// 	<usage> 	 <brief> (vertical alignment)
	//    ...
	ta.SetHead("Commands:")

	for _, command := range p.commands {
		for _, usage := range command.Usages {
			ta.AddRow(usage, command.Brief)
		}
	}

	ta.AlignColumn(0)

	lines, _ := ta.Build(trav.GetConfig().GetTabSize(), true)

	err = trav.AddLines(lines)
	if err != nil {
		return err
	}

	// Flags:
	// 	<flags>
	if p.flags != nil {
		ta.SetHead("Flags:")

		for _, flag := range p.flags {
			str := strings.Join(flag.Usages, settings.Spacing)

			ta.AddRow(str, flag.Brief)
		}

		ta.AlignColumn(0)

		lines, _ := ta.Build(tabSize, true)

		err = trav.AddLines(lines)
		if err != nil {
			return err
		}
	}

	return nil
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
		_, ok := p.commands[cmd.Name]
		if !ok {
			p.commands[cmd.Name] = cmd
		}
	}
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
	msg := pd.NewInputMsg(text)

	err := p.display.Send(msg)
	if err != nil {
		return "", err
	}

	resp, err := msg.Receive()
	if err != nil {
		return "", err
	}

	return resp, nil
}

// Logf logs a formatted line to the program's buffer.
//
// Parameters:
//   - format: The format of the line.
//   - args: The items to log.
//
// Returns:
//   - error: An error if the program stopped abruptly
//     (due to call to Program.Panic())
func (p *Program) Logf(format string, args ...interface{}) error {
	str := fmt.Sprintf(format, args...)

	msg := pd.NewLogMsg(str)

	err := p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Logln logs a line to the program's buffer.
//
// Parameters:
//   - args: The items to log.
//
// Returns:
//   - error: An error if the program stopped abruptly
//     (due to call to Program.Panic())
func (p *Program) Logln(args ...interface{}) error {
	str := fmt.Sprintln(args...)
	str = strings.TrimSuffix(str, "\n")

	msg := pd.NewLogMsg(str)

	err := p.display.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
