package simple

import (
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// Program is a struct that represents a program.
type Program struct {
	// Name is the name of the program.
	Name string

	// Version is the version of the program.
	Version string

	// command_table is the table of commands.
	command_table map[string]*Command
}

// Fix implements the errors.Fixer interface.
func (p *Program) Fix() error {
	if p == nil {
		return gcers.NilReceiver
	}

	name := strings.TrimSpace(p.Name)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	p.Name = name

	p.Version = strings.TrimSpace(p.Version)

	if p.command_table == nil {
		p.command_table = make(map[string]*Command)
	} else {
		for k, cmd := range p.command_table {
			err := gcers.Fix("command "+strconv.Quote(k), cmd, false)
			if err != nil {
				return err
			}
		}
	}

	// Add version command if needed.
	if p.Version != "" {
		ok := p.HasCommand("version")
		if !ok {
			version_cmd := &Command{
				Name: "version",
				RunFn: func(p *Program, _ []string) error {
					_, err := fmt.Println(p.Version)
					return err
				},
				Argument: NoArguments,
			}

			p.command_table["version"] = version_cmd
		}
	}

	return nil
}

// AddCommands adds commands to the program.
//
// Parameters:
//   - commands: The commands to add.
//
// Nil commands are ignored.
func (p *Program) AddCommands(commands ...*Command) {
	if p == nil || len(commands) == 0 {
		return
	}

	if p.command_table == nil {
		p.command_table = make(map[string]*Command)
	}

	for _, command := range commands {
		if command == nil {
			continue
		}

		p.command_table[command.Name] = command
	}
}

// HasCommand checks if the program has a command with the given name.
//
// Returns:
//   - bool: True if the program has a command with the given name, false otherwise.
func (p Program) HasCommand(name string) bool {
	if p.command_table == nil {
		return false
	}

	_, ok := p.command_table[name]
	return ok
}

// RetrieveCommand retrieves the command with the given name.
//
// Returns:
//   - *Command: The command with the given name. Never returns nil.
//   - bool: True if the program has a command with the given name, false otherwise.
func (p Program) RetrieveCommand(name string) (*Command, bool) {
	if p.command_table == nil {
		return nil, false
	}

	cmd, ok := p.command_table[name]
	if !ok {
		return nil, false
	}

	return cmd, true
}

// Command is a method that returns an iterator of commands.
//
// Returns:
//   - iter.Seq2[string, *Command]: The iterator of commands.
func (p Program) Command() iter.Seq2[string, *Command] {
	return func(yield func(string, *Command) bool) {
		for k, cmd := range p.command_table {
			if !yield(k, cmd) {
				break
			}
		}
	}
}

// Run is a method that runs the program.
//
// Parameters:
//   - args: The arguments to run the program with. This is os.Args.
//
// Returns:
//   - error: The error that occurred. Never returns nil.
func (p Program) Run(args []string) error {
	if len(args) < 2 {
		_, err := fmt.Println("Usage:", p.Name, "<cmd>")
		if err != nil {
			return err
		}

		_, err = fmt.Println("Use \"help\" command to see the list of available commands")
		if err != nil {
			return err
		}

		return nil
	}

	command := args[1]

	cmd, ok := p.command_table[command]
	if !ok {
		_, err := fmt.Println("Unknown command: " + command)
		if err != nil {
			return err
		}

		_, err = fmt.Println("Use \"help\" command to see the list of available commands")
		if err != nil {
			return err
		}

		return nil
	}

	args, err := cmd.parse(args[2:])
	if err != nil {
		return err
	}

	err = cmd.RunFn(&p, args)
	if err != nil {
		return fmt.Errorf("command %q failed: %w", command, err)
	}

	return nil
}

// Print is a method that prints the given arguments. A newline is added at the end.
//
// Parameters:
//   - args: The arguments to print.
//
// Returns:
//   - error: The error that occurred.
func (p Program) Print(args ...any) error {
	_, err := fmt.Println(args...)
	return err
}

// Printf is a method that prints the given format and arguments. A newline is added at the end.
//
// Parameters:
//   - format: The format to print.
//   - args: The arguments to print.
//
// Returns:
//   - error: The error that occurred.
func (p Program) Printf(format string, args ...any) error {
	_, err := fmt.Printf(format+"\n", args...)
	return err
}

// PrintNewline is a method that prints a newline.
//
// Returns:
//   - error: The error that occurred.
func (p Program) PrintNewline() error {
	_, err := fmt.Println()
	return err
}

// DefaultExitSequence is a method that prints the given error and exits the program.
//
// Parameters:
//   - err: The error to print.
func DefaultExitSequence(err error) {
	var exit_code int

	if err == nil {
		_, err := fmt.Println("Command ran successfully")
		if err != nil {
			panic(err)
		}

		exit_code = 0
	} else {
		_, err := fmt.Println(err.Error())
		if err != nil {
			panic(err)
		}

		exit_code = 1
	}

	_, err = fmt.Println()
	if err != nil {
		panic(err)
	}

	_, err = fmt.Println("Press ENTER to exit...")
	if err != nil {
		panic(err)
	}

	_, err = fmt.Scanln()
	if err != nil {
		panic(err)
	}

	os.Exit(exit_code)
}
