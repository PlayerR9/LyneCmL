package Complex

import (
	"fmt"
	"strings"

	util "github.com/PlayerR9/LyneCmL/Complex/util"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// RunFunc is a function that will be executed when the command is called.
//
// Parameters:
//   - p: The program that the command is being executed on.
//   - args: The arguments that were passed to the command.
//   - data: The data that was passed to the command. (if any)
//
// Returns:
//   - any: The result of the command.
//   - error: An error if the command failed to execute.
type RunFunc func(p *Program, args []string, data any) (any, error)

var (
	// NoRunFunc is a function that does nothing.
	NoRunFunc RunFunc
)

func init() {
	NoRunFunc = func(p *Program, args []string, data any) (any, error) {
		return data, nil
	}
}

// Command is a command that a program can execute.
type Command struct {
	// Name is the name of the command.
	Name string

	// Usage is the usage of the command.
	Usage string

	// Brief is a brief description of the command.
	Brief string

	// Description is a description of the command.
	Description *Description

	// Argument is the argument of the command.
	Argument *Argument

	// Run is the function that will be executed when the command is called.
	Run RunFunc

	// subCommands is a map of sub-commands that the command can execute.
	// If at least one sub-command is present, then the first argument
	// of the command will be the sub-command.
	//
	// Then the return value of the run function will be passed to the
	// sub-command as the data argument.
	subCommands map[string]*Command
}

// fix fixes the command by trimming all the strings and setting default values.
func (c *Command) fix() {
	c.Name = strings.TrimSpace(c.Name)
	c.Usage = strings.TrimSpace(c.Usage)
	c.Brief = strings.TrimSpace(c.Brief)

	if c.Argument == nil {
		c.Argument = NoArgument
	}

	c.Argument.fix()

	if c.Run == nil {
		c.Run = NoRunFunc
	}
}

// FString implements the ffs.FStringer interface.
func (c *Command) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	// Command: <name>
	err := trav.AddJoinedLine(" ", "Command:", c.Name)
	if err != nil {
		return err
	}

	trav.EmptyLine()

	// Usage: <usage>
	err = trav.AddJoinedLine(" ", "Usage:", c.Usage)
	if err != nil {
		return err
	}

	if c.Description == nil {
		return nil
	}

	// Description:
	// 	<description>
	trav.EmptyLine()

	err = trav.AddLine("Description:")
	if err != nil {
		return err
	}

	err = ffs.ApplyForm(
		trav.GetConfig(
			ffs.WithModifiedIndent(1),
		),
		trav,
		c.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

// AddSubCommand adds a sub-command to the command.
//
// Parameters:
//   - cmds: The sub-commands to add.
//
// Behaviors:
//   - If a sub-command has the same name as another sub-command, the first one
//     will be kept.
//   - nil sub-commands will be filtered out.
func (c *Command) AddSubCommand(cmds ...*Command) {
	cmds = us.FilterNilValues(cmds)

	if len(cmds) == 0 {
		return
	}

	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}

	for _, cmd := range cmds {
		cmd.fix()
		if cmd.Name == "" {
			continue
		}

		_, ok := c.subCommands[cmd.Name]
		if !ok {
			c.subCommands[cmd.Name] = cmd
		}
	}
}

// Parsed is a parsed command.
type Parsed struct {
	// args are the arguments that were parsed.
	args []string

	// data is the data that was parsed. (if any)
	data any
}

// parseArgs parses the arguments and runs the command.
//
// Parameters:
//   - p: The program that the command is being executed on.
//   - args: The arguments that were passed to the command.
//   - cmd: The command to run.
//
// Returns:
//   - *Parsed: The parsed arguments and data.
//   - error: An error if the command failed to execute.
func parseArgs(p *Program, args []string, cmd *Command) (*Parsed, error) {
	if len(cmd.subCommands) == 0 || len(args) == 0 {
		parsed, err := handleCmd(p, args, cmd)
		return parsed, err
	}

	// Recursive case
	subCmd, ok := cmd.subCommands[args[0]]
	if !ok {
		return nil, util.NewErrUnknownCommand(args[0])
	}

	parsed, err := parseArgs(p, args[1:], subCmd)
	if err != nil {
		return nil, fmt.Errorf("in sub-command %q: %w", subCmd.Name, err)
	}

	// Handle the command
	parsed, err = handleCmd(p, parsed.args, subCmd)
	return parsed, err
}

// handleCmd handles the command by validating the arguments and running the command.
//
// Parameters:
//   - p: The program that the command is being executed on.
//   - args: The arguments that were passed to the command.
//   - cmd: The command to run.
//
// Returns:
//   - *Parsed: The parsed arguments and data.
//   - error: An error if the command failed to execute.
func handleCmd(p *Program, args []string, cmd *Command) (*Parsed, error) {
	validatedArgs, err := cmd.Argument.validate(args)
	if err != nil {
		return nil, err
	}

	data, n, err := cmd.Argument.parseFunc(p, validatedArgs)
	if err != nil {
		return nil, err
	}

	if n < 0 {
		n = 0
	} else if n > len(args) {
		n = len(args)
	}

	res, err := cmd.Run(p, validatedArgs[:n], data)
	if err != nil {
		return nil, fmt.Errorf("error running command: %w", err)
	}

	return &Parsed{
		args: args[n:],
		data: res,
	}, nil
}
