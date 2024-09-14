package simple

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ds "github.com/PlayerR9/LyneCml/screen"
	"github.com/PlayerR9/LyneCml/simple/internal"
	cmls "github.com/PlayerR9/LyneCml/style"
	fs "github.com/PlayerR9/go-commons/Formatting/strings"
	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/gdamore/tcell"
)

// Program is a struct that represents a program.
type Program struct {
	// FullName is the full name of the program.
	FullName string

	// Name is the name of the program. This is used when typing commands.
	Name string

	// Version is the version of the program. Leave empty if not needed.
	Version string

	// Brief is a brief description of the program. Leave empty if not needed.
	Brief string

	// Description is the description of the program. Leave empty if not needed.
	Description []string

	// command_list is the list of commands.
	command_list map[string]*Command

	// screen is the screen of the program.
	screen *ds.Screen

	// style is the style of the program.
	style cmls.Style[cmls.ColorType]
}

// HelpLines is a method that returns the help lines of the program.
//
// Returns:
//   - []string: The help lines of the program.
func (p Program) HelpLines() []string {
	var lines []string

	if p.Brief != "" {
		lines = append(lines, p.FullName+" — "+p.Brief)
	} else {
		lines = append(lines, p.FullName)
	}

	lines = append(lines, "")

	if len(p.Description) > 0 {
		lines = append(lines, p.Description...)
		lines = append(lines, "", "")
	}

	lines = append(lines, "Usage:")

	var all_commands [][]string

	for _, cmd := range p.command_list {
		arr := cmd.Usage()
		arr[0] = p.Name + " " + arr[0]

		all_commands = append(all_commands, []string{arr[0], arr[1]})
	}

	table, err := fs.TableEntriesAlign(all_commands, 3)
	if err != nil {
		panic(fmt.Sprintf("failed to create table: %v", err))
	}

	for _, row := range table {
		lines = append(lines, strings.Join(row, ""))
	}

	return nil
}

// Fix is a method that builds and fixes the program. Remember to call
// this method before running the program.
//
// Returns:
//   - error: An error of if the program could not be fixed.
func (p *Program) Fix() error {
	if p == nil {
		return gcers.NilReceiver
	}

	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return errors.New("program name cannot be empty")
	}

	if p.command_list == nil {
		p.command_list = make(map[string]*Command)
	} else {
		for k, cmd := range p.command_list {
			err := cmd.Fix()
			if err != nil {
				return fmt.Errorf("failed to fix command %q: %w", k, err)
			}
		}
	}

	help_cmd := &Command{
		Name:  "help",
		Brief: "Displays the help message.",
		Description: NewDescription(
			"The help command displays the help message for the program or for a specific command.",
			"If no command is specified, the help command will display the help message for the program.",
			"The help command is useful for getting help on the program or on a specific command.",
		).
			Build(),
		RunFunc: func(ctx context.Context, args []string) error {
			if len(args) == 0 {
				lines := p.HelpLines()

				for i, line := range lines {
					err := Do(ctx, NewActPrint(tcell.StyleDefault, line))
					if err != nil {
						return gcint.
					}

					err := p.Print(line)
					if err != nil {
						return err
					}
				}
			} else {
				name := args[0]

				cmd, ok := p.command_list[name]
				if !ok {
					return fmt.Errorf("command %q not found", name)
				}

				lines := cmd.HelpLines()

				for _, line := range lines {
					err := p.Print(line)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
		Argument: AtMostNArgs("command", 1),
	}

	err := help_cmd.Fix()
	if err != nil {
		panic(fmt.Sprintf("failed to fix help command: %v", err))
	}

	p.command_list["help"] = help_cmd

	if p.Version != "" {
		version_cmd := &Command{
			Name: "version",
			RunFunc: func(ctx context.Context, _ []string) error {
				err := p.Print(p.Version)
				if err != nil {
					return err
				}

				return nil
			},
			Argument: NoArguments,
		}

		err := version_cmd.Fix()
		if err != nil {
			panic(fmt.Sprintf("failed to fix version command: %v", err))
		}

		p.command_list["version"] = version_cmd
	}

	return nil
}

// AddCommand adds a command to the program. Ignores nil commands.
//
// Parameters:
//   - cmd: The command to add.
func (p *Program) AddCommand(cmd *Command) {
	if p == nil || cmd == nil {
		return
	}

	if p.command_list == nil {
		p.command_list = make(map[string]*Command)
	}

	p.command_list[cmd.Name] = cmd
}

// AddCommands is a convenience method that adds multiple commands to the
// program. It is the same as calling AddCommand for each command.
//
// Parameters:
//   - cmds: The commands to add.
//
// Returns:
//   - error: An error if the command failed to fix.
func (p *Program) AddCommands(cmds ...*Command) {
	if p == nil {
		return
	}

	var top int

	for i := 0; i < len(cmds); i++ {
		if cmds[i] != nil {
			cmds[top] = cmds[i]
			top++
		}
	}

	cmds = cmds[:top:top]

	if len(cmds) == 0 {
		return
	}

	if p.command_list == nil {
		p.command_list = make(map[string]*Command)
	}

	for _, cmd := range cmds {
		p.command_list[cmd.Name] = cmd
	}
}

// Run is a method that runs the program.
//
// Parameters:
//   - args: The arguments passed to the program. It should be os.Args.
//
// Returns:
//   - error: An error of type *errors.Err[ErrorCode] if there was an error.
func (p *Program) Run(bg_style tcell.Style, args []string) error {
	ctx := ds.NewContext(ds.WithDefaultStyle(bg_style))

	err := ds.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	if len(args) < 2 {
		return internal.NewErrMissingCommand()
	}

	command := args[1]

	if p.command_list == nil {
		return gcers.NewErrInvalidUsage(errors.New("program is in an invalid state"), "Please call Fix() before calling Run()")
	}

	cmd, ok := p.command_list[command]
	if !ok {
		return internal.NewErrInvalidCommand(command)
	}

	screen, err := ds.NewScreen(bg_style)
	if err != nil {
		return fmt.Errorf("failed to create screen: %w", err)
	}

	p.screen = screen

	args = args[2:]

	parsed, err := cmd.parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}

	// args_left := args[len(parsed):]

	err = cmd.RunFunc(p, parsed)
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}

// SetCell is a method that sets the cell at the given x and y coordinates on the screen to the given character
// with the given style.
//
// Parameters:
//   - x: The x-coordinate of the cell to set.
//   - y: The y-coordinate of the cell to set.
//   - char: The character to set the cell to.
//   - style: The style to set the cell to.
func (p Program) DrawCell(x, y int, char rune, style tcell.Style) {

}

/* // BgStyle returns the background style.
//
// Returns:
//   - tcell.Style: The background style.
func (p Program) BgStyle() tcell.Style {
	return ds.LightModeStyle
} */

// Height returns the height of the screen.
//
// Returns:
//   - int: The height of the screen.
func (p Program) Height() int {
	return p.screen.Height()
}
