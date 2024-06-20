package Simple

import (
	"fmt"
	"slices"
	"strings"

	util "github.com/PlayerR9/LyneCmL/Simple/util"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// parseFlags is a function that parses flags.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - error: An error if the flags failed to parse.
//
// Format:
//
//	--flag -> for boolean flags
//	--flag=arg -> for flags with arguments
func parseFlags(flagMap map[string]any, args []string) error {
	var flagList []any

	for _, value := range flagMap {
		flagList = append(flagList, value)
	}

	for len(args) > 0 {
		fields := strings.SplitN(args[0], "=", 2)

		switch len(fields) {
		case 1:
			// Boolean flag
			left := fields[0]

			f := func(flag any) bool {
				bf, ok := flag.(*BoolFlag)
				if !ok {
					return false
				}

				return bf.Name == left
			}

			index := slices.IndexFunc(flagList, f)

			if index == -1 {
				return fmt.Errorf("unknown flag %q", left)
			}

			flag := flagList[index].(*BoolFlag)
			flag.value = true

			flagList = slices.Delete(flagList, index, index+1)
		case 2:
			left := fields[0]
			right := fields[1]

			f := func(flag any) bool {
				valFlag, ok := flag.(*Flag[any])
				if !ok {
					return false
				}

				return valFlag.Name == left
			}

			index := slices.IndexFunc(flagList, f)
			if index == -1 {
				return fmt.Errorf("unknown flag %q", left)
			}

			flag := flagList[index].(*Flag[any])

			value, err := flag.ParseFunc(right)
			if err != nil {
				return fmt.Errorf("in flag %q: %w", flag.Name, err)
			}

			flag.SetValue(value)

			flagList = slices.Delete(flagList, index, index+1)
		default:
			return fmt.Errorf("empty argument")
		}
	}

	for _, flag := range flagList {
		switch flag := flag.(type) {
		case *Flag[any]:
			flag.SetValue(flag.DefaultVal)
		case *BoolFlag:
			flag.SetValue(false)
		}
	}

	return nil
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

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		err = ue.NewErrPanic(r)

		err := p.Panic(err)
		if err != nil {
			panic(err)
		}
	}()

	data, n, err := cmd.Argument.parseFunc(p, validatedArgs)
	if err != nil {
		return nil, err
	}

	ok := p.display.IsDone()
	if ok {
		return nil, nil
	}

	if n < 0 {
		n = 0
	} else if n > len(args) {
		n = len(args)
	}

	err = parseFlags(cmd.flags, args[n:])
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
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
