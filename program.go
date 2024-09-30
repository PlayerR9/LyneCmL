package cml

import (
	"fmt"
	"iter"
	"os"
	"strings"
	"text/tabwriter"

	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-errors/assert"
	"github.com/PlayerR9/go-sets"
)

// Program is a collection of commands
type Program struct {
	// Name is the name of the program
	Name string

	// Version is the version of the program
	Version string

	// Brief is the brief description of the program
	Brief string

	// commands is a map of command names to command objects
	commands *sets.OrderedSet[string, *Command]
}

// Write implements the io.Writer interface.
func (p Program) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

// Fix implements the errors.Fixer interface.
func (p *Program) Fix() error {
	if p == nil {
		err := gers.NewErrNilReceiver()
		err.AddFrame("Fix()")

		return err
	}

	p.Name = strings.TrimSpace(p.Name)
	p.Brief = strings.TrimSpace(p.Brief)

	if p.commands == nil {
		p.commands = sets.NewOrderedSet[string, *Command]()
	} else {
		fixed_commands := sets.NewOrderedSet[string, *Command]()

		for _, cmd := range p.commands.Entry() {
			if cmd == nil {
				continue
			}

			err := cmd.Fix()
			if err != nil {
				return err
			}

			fixed_commands.Add(cmd.Name, cmd)
		}

		p.commands = fixed_commands
	}

	p.Version = strings.TrimSpace(p.Version)

	if p.Version != "" {
		ok := p.commands.Contains("version")
		if !ok {
			version_cmd := &Command{
				Name:  "version",
				Brief: "Show the version of the program",
				RunFn: func(p *Program, _ []string) error {
					assert.NotNil(p, "p")

					return p.Print("Version:", p.Version)
				},
			}

			p.commands.ForceAdd("version", version_cmd)
		}
	}

	ok := p.commands.Contains("help")
	if !ok {
		help_cmd := &Command{
			Name:  "help",
			Brief: "Show this help",
			RunFn: func(p *Program, _ []string) error {
				err := p.Print("Available commands:")
				if err != nil {
					return err
				}

				w := p.TabWriter(3, 3, 1)
				defer w.Flush()

				for name, cmd := range p.Command() {
					err := w.PrintRow("\t", name, cmd.Brief)
					if err != nil {
						return err
					}
				}

				return nil
			},
		}

		p.commands.ForceAdd("help", help_cmd)
	}

	return nil
}

// AddCommand adds a command if it is not nil.
//
// Parameters:
//   - cmd: The command to add.
//
// Does nothing if the receiver is nil.
func (p *Program) AddCommand(cmd *Command) {
	if p == nil || cmd == nil {
		return
	}

	if p.commands == nil {
		p.commands = sets.NewOrderedSet[string, *Command]()
	}

	p.commands.ForceAdd(cmd.Name, cmd)
}

// AddCommands is a more convenient way to add multiple commands.
//
// Parameters:
//   - cmds: The commands to add.
//
// Does nothing if the receiver is nil.
func (p *Program) AddCommands(cmds ...*Command) {
	if p == nil {
		return
	}

	var top int

	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]

		if cmd != nil {
			cmds[top] = cmd
			top++
		}
	}

	cmds = cmds[:top:top]
	if len(cmds) == 0 {
		return
	}

	if p.commands == nil {
		p.commands = sets.NewOrderedSet[string, *Command]()
	}

	for _, cmd := range cmds {
		p.commands.ForceAdd(cmd.Name, cmd)
	}
}

// Run runs the program.
//
// Parameters:
//   - args: The command line arguments. This is expected to be os.Args.
//
// Returns:
//   - error: Any error that may have occurred.
func (p *Program) Run(args []string) error {
	if p == nil {
		return nil
	}

	if len(args) < 2 {
		fmt.Println("Use \"help\" for a list of commands")

		return fmt.Errorf("no command specified")
	}

	command := args[1]

	if p.Name == "" {
		p.Name = args[0]
	}

	if p.commands == nil {
		return fmt.Errorf("no commands")
	}

	cmd, ok := p.commands.Get(command)
	if !ok {
		fmt.Println("Use \"help\" for a list of commands")

		return fmt.Errorf("unknown command: %s", command)
	}

	args, err := cmd.parse_args(args[2:])
	if err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	err = cmd.RunFn(p, args)
	return err
}

// Print prints the given arguments. A newline is always printed at the end
// and a space is added between arguments.
//
// No arguments will print a newline.
//
// Parameters:
//   - args: The arguments to print.
//
// Returns:
//   - error: Any error that may have occurred.
func (p Program) Print(args ...any) error {
	var err error

	if len(args) == 0 {
		_, err = fmt.Println()
	} else {
		_, err = fmt.Println(args...)
	}

	return err
}

// Printf prints the given format and arguments. However,
// a newline is always printed at the end.
//
// Parameters:
//   - format: The format to print.
//   - args: The arguments to print.
//
// Returns:
//   - error: Any error that may have occurred.
func (p Program) Printf(format string, args ...any) error {
	_, err := fmt.Printf(format, args...)
	if err != nil {
		return err
	}

	_, err = fmt.Println()
	return err
}

// TabWriter returns a TabWriter. For the arguments, see tabwriter.NewWriter.
//
// Parameters:
//   - min_width: The minimum width.
//   - tab_width: The tab width.
//   - padding: The padding.
//
// Returns:
//   - TabWriter: The TabWriter.
func (p *Program) TabWriter(min_width, tab_width, padding int) TabWriter {
	w := tabwriter.NewWriter(p, min_width, tab_width, padding, ' ', 0)

	return TabWriter{
		w: w,
	}
}

// Command returns a sequence of commands.
//
// Parameters:
//   - name: The name of the command.
//
// Returns:
//   - iter.Seq2[string, *Command]: A sequence of commands. Never returns nil.
func (p Program) Command() iter.Seq2[string, *Command] {
	if p.commands == nil {
		return func(yield func(string, *Command) bool) {}
	}

	return p.commands.Entry()
}

// DefaultExitSequence is the default exit sequence.
//
// In the default exit sequence, the error is printed if it is not nil, and
// "Success" is printed if it is nil. The user is then prompted to press ENTER
// to exit.
//
// Parameters:
//   - err: The error that may have occurred.
func DefaultExitSequence(err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}

	fmt.Println()

	fmt.Println("Press ENTER to exit...")

	fmt.Scanln()
}
