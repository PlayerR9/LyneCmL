package simple

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/PlayerR9/LyneCml/simple/internal"
	gcers "github.com/PlayerR9/go-commons/errors"
	dbg "github.com/PlayerR9/go-debug/assert"
)

// Program is a struct that represents a program.
type Program struct {
	// FullName is the full name of the program.
	FullName string

	// Name is the name of the program. This is used when typing commands.
	Name string

	// Version is the version of the program. Leave empty if not needed.
	Version string

	// command_list is the list of commands.
	command_list map[string]*Command
}

// Fix is a method that builds and fixes the program. Remember to call
// this method before running the program.
//
// Returns:
//   - error: An error of type *errors.Err[ErrorCode] if there was an error.
func (p *Program) Fix() error {
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
		Name: "help",
		RunFunc: func(p *Program, _ []string) error {
			var lines []string

			lines = append(lines, p.FullName)
			lines = append(lines, "")
			lines = append(lines, "Usage:")

			for _, cmd := range p.command_list {
				lines = append(lines, fmt.Sprintf("%s %s", p.Name, cmd.Usage()))
			}

			err := p.Print(strings.Join(lines, "\n"))
			if err != nil {
				return err
			}

			return nil
		},
		Argument: NoArguments,
	}

	err := help_cmd.Fix()
	dbg.AssertErr(err, "help_cmd.Fix()")

	p.command_list["help"] = help_cmd

	if p.Version != "" {
		version_cmd := &Command{
			Name: "version",
			RunFunc: func(p *Program, _ []string) error {
				err := p.Print(p.Version)
				if err != nil {
					return err
				}

				return nil
			},
			Argument: NoArguments,
		}

		err := version_cmd.Fix()
		dbg.AssertErr(err, "version_cmd.Fix()")

		p.command_list["version"] = version_cmd
	}

	return nil
}

// AddCommand adds a command to the program. Ignores nil commands.
//
// Parameters:
//   - cmd: The command to add.
func (p *Program) AddCommand(cmd *Command) {
	if cmd == nil {
		return
	}

	if p.command_list == nil {
		p.command_list = make(map[string]*Command)
	}

	p.command_list[cmd.Name] = cmd
}

// Run is a method that runs the program.
//
// Parameters:
//   - args: The arguments passed to the program. It should be os.Args.
//
// Returns:
//   - error: An error of type *errors.Err[ErrorCode] if there was an error.
func (p Program) Run(args []string) error {
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

	args = args[2:]

	parsed, err := cmd.parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}

	// args_left := args[len(parsed):]

	err = cmd.RunFunc(&p, parsed)
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}

// Print is a method that prints a message to the standard output with a newline.
//
// Parameters:
//   - a: The arguments to print.
//
// Returns:
//   - error: An error if the message could not be printed.
func (p Program) Print(a ...any) error {
	_, err := fmt.Println(a...)
	if err != nil {
		return err
	}

	return nil
}

// Printf is a method that prints a formatted message to the standard output with a newline.
//
// Parameters:
//   - format: The format of the message.
//   - a: The arguments to print.
//
// Returns:
//   - error: An error if the message could not be printed.
func (p Program) Printf(format string, a ...any) error {
	a = append(a, "\n")

	_, err := fmt.Printf(format, a...)
	if err != nil {
		return err
	}

	return nil
}

// DefaultExitSequence is the default exit sequence for the program.
//
// Parameters:
//   - err: The error that occurred. If nil, the program will exit with code 0.
func DefaultExitSequence(err error) {
	var exit_code int

	if err == nil {
		fmt.Println("Done!")
		exit_code = 0
	} else {
		fmt.Println(err.Error())

		switch err := err.(type) {
		case *gcers.Err[internal.ErrorCode]:
			fmt.Println()
			fmt.Println("Suggestions:")

			for _, suggestion := range err.Suggestions {
				fmt.Println("\t", suggestion)
			}

			exit_code = int(err.Code) + 2
		default:
			exit_code = 1
		}
	}

	fmt.Println()
	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()

	os.Exit(exit_code)
}
