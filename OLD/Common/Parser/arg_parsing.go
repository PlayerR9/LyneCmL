package Parser

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	cms "github.com/PlayerR9/LyneCml/OLD/Simple"
	gcers "github.com/PlayerR9/go-commons/errors"
	luch "github.com/PlayerR9/go-commons/runes"
	lls "github.com/PlayerR9/go-commons/stack"
)

const (
	// LongFlagPrefix is the prefix for long flags.
	LongFlagPrefix string = "--"

	// ShortFlagPrefix is the prefix for short flags.
	ShortFlagPrefix string = "-"
)

// flagParser is a parser for flags.
type flagParser struct {
	// flagList is the list of flags to parse.
	flagList []*cms.Flag

	// args is the list of arguments to parse.
	args []string
}

// findFlag finds the index of the flag.
//
// Parameters:
//   - name: The name of the flag.
//   - isShort: True if the flag is a short flag. False otherwise.
//
// Returns:
//   - int: The index of the long flag. -1 if not found.
func (p *flagParser) findFlag(name string, isShort bool) int {
	var search func(flag *cms.Flag) bool

	if isShort {
		r, _ := utf8.DecodeRuneInString(name)

		search = func(flag *cms.Flag) bool {
			return flag.ShortName == r
		}
	} else {
		search = func(flag *cms.Flag) bool {
			return flag.LongName == name
		}
	}

	index := slices.IndexFunc(p.flagList, search)

	return index
}

// parseSingleFlag parses a single flag.
//
// Format:
//
//	--flag		: only boolean flags without arguments
//	--flag=value	: any flag that requires an argument
//	--flag value 	: any flag that requires an argument
//
//	-flag		: only boolean flags without arguments
//	-flag=value	: any flag that requires an argument
//	-flag value 	: any flag that requires an argument
//	-abc		: multiple short flags
//
// Parameters:
//   - header: The header of the flag.
//   - right: The right side of the flag.
//   - isShort: True if the flag is a short flag. False otherwise.
//
// Returns:
//   - string: The remaining right side of the flag.
//   - error: An error if the flag failed to parse.
func (fp *flagParser) parseSingleFlag(header string, right string, isShort bool) (string, error) {
	index := fp.findFlag(header, isShort)
	if index == -1 {
		return "", NewErrFlagNotFound(isShort, header)
	}

	flag := fp.flagList[index]

	if flag.Argument == nil {
		if !isShort && right != "" {
			return "", fmt.Errorf("flag --%s does not take an argument", flag.LongName)
		}

		flag.Apply("")

		fp.flagList = slices.Delete(fp.flagList, index, index+1)

		if isShort {
			return right, nil
		}
	}

	var todel int

	if right == "" {
		if len(fp.args) == 1 {
			return "", NewErrFlagMissingArg(isShort, flag)
		}

		right = fp.args[1]

		todel = 2
	} else {
		todel = 1
	}

	err := flag.Apply(right)
	if err != nil {
		return "", NewErrInvalidFlag(isShort, flag, err)
	}

	fp.args = fp.args[todel:]

	fp.flagList = slices.Delete(fp.flagList, index, index+1)

	return "", nil
}

// isValidFlag checks if the flag is valid.
//
// Parameters:
//   - arg: The argument to check.
//
// Returns:
//   - bool: True if the flag is short. False otherwise.
//   - bool: True if the flag is valid. False otherwise.
func isValidFlag(arg string) (bool, bool) {
	ok := strings.HasPrefix(arg, LongFlagPrefix)
	if ok {
		return false, true
	}

	ok = strings.HasPrefix(arg, ShortFlagPrefix)

	return true, ok
}

// parseOne parses one flag.
//
// Parameters:
//   - header: The header of the flag.
//   - isShort: True if the flag is a short flag. False otherwise.
//
// Returns:
//   - error: An error if the flag failed to parse.
func (fp *flagParser) parseOne(header string, isShort bool) error {
	fields := strings.SplitN(header, "=", 2)
	header = fields[0]

	var right string

	if len(fields) == 2 {
		right = fields[1]
	} else {
		right = ""
	}

	if !isShort {
		_, err := fp.parseSingleFlag(header, right, isShort)
		if err != nil {
			return err
		}

		return nil
	}

	runes, err := luch.StringToUtf8(header)
	if err != nil {
		return err
	}

	if len(runes) == 1 {
		right, err = fp.parseSingleFlag(header, right, true)
		if err != nil {
			ok := gcers.Is[*ErrFlagNotFound](err)
			if ok {
				err = nil
			}
		} else if right != "" {
			err = fmt.Errorf("extra argument %q", right)
		}
	} else {
		for _, letter := range runes {
			right, err = fp.parseSingleFlag(string(letter), right, true)
			if err != nil {
				ok := gcers.Is[*ErrFlagNotFound](err)
				if ok {
					err = fmt.Errorf("invalid merged short flags -%s", header)
				}

				return err
			}
		}

		if right != "" {
			err = fmt.Errorf("extra argument %q", right)
		}
	}

	return err
}

type parsingResult struct {
	cutSet    []string
	data      any
	argsLeft  []string
	flagsLeft []*cms.Flag
}

type flagResult struct {
	flagList []*cms.Flag
	args     []string
}

func parseFlags(flagMap map[string]*cms.Flag, args []string) (*flagResult, error) {
	var flagList []*cms.Flag

	for _, value := range flagMap {
		flagList = append(flagList, value)
	}

	fp := &flagParser{
		flagList: flagList,
		args:     args,
	}

	for len(fp.args) > 0 {
		header := fp.args[0]

		isShort, ok := isValidFlag(header)
		if !ok {
			fr := &flagResult{
				flagList: fp.flagList,
				args:     fp.args,
			}

			return fr, nil
		}

		if isShort {
			header = strings.TrimPrefix(header, ShortFlagPrefix)
		} else {
			header = strings.TrimPrefix(header, LongFlagPrefix)
		}

		err := fp.parseOne(header, isShort)
		if err != nil {
			return nil, NewErrInvalidArg(header, isShort, err)
		}

		if len(fp.flagList) == 0 {
			fr := &flagResult{
				flagList: nil,
				args:     fp.args,
			}
			return fr, nil
		}
	}

	fr := &flagResult{
		flagList: fp.flagList,
		args:     nil,
	}

	return fr, nil
}

// parseArguments parses the arguments.
//
// Parameters:
//   - a: The argument to parse.
//   - flagMap: The map of flags.
//   - args: The arguments to parse.
//
// Returns:
//   - *parsingResult: The parsed arguments and data.
//   - error: An error if the arguments failed to parse.
func parseArguments(a *cms.Argument, flagMap map[string]*cms.Flag, args []string) (*parsingResult, error) {
	index := -1

	for i := 0; i < len(args); i++ {
		ok := strings.HasPrefix(args[i], "-")
		if ok {
			index = i
			break
		}
	}

	var actualArgs []string

	if index == -1 {
		actualArgs = args
	} else {
		actualArgs = args[:index]
	}

	parsed, err := a.Apply(actualArgs)
	if err != nil {
		return nil, err
	}

	fr, err := parseFlags(flagMap, args[parsed.Idx:])
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	result := &parsingResult{
		cutSet:    parsed.CutSet,
		data:      parsed.Data,
		argsLeft:  fr.args,
		flagsLeft: fr.flagList,
	}

	return result, nil
}

// ParseArgs parses the arguments.
//
// Parameters:
//   - commandMap: The map of commands.
//   - args: The arguments to parse.
//
// Returns:
//   - lls.Stacker[*cms.ExecProcess]: The list of commands to execute.
//   - error: An error if the arguments failed to parse.
func ParseArgs(commandMap map[string]*cms.Command, args []string) (lls.Stacker[*cms.ExecProcess], error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments")
	}

	p := &Parser{
		argsLeft: args,
		execList: lls.NewStack[*cms.ExecProcess](),
	}

	cmd, ok := commandMap[args[0]]
	if !ok {
		return nil, fmt.Errorf("command %q not found", args[0])
	}

	p.argsLeft = p.argsLeft[1:]

	err := p.parseArgs(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if len(p.argsLeft) != 0 {
		return nil, fmt.Errorf("extra arguments %q", p.argsLeft)
	}

	return p.execList, nil
}

type Parser struct {
	argsLeft []string

	execList *lls.Stack[*cms.ExecProcess]
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
func (p *Parser) handleCmd(cmd *cms.Command) error {
	fm := cmd.GetFlagMap()

	parsed, err := parseArguments(cmd.Argument, fm, p.argsLeft)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	elem := cms.NewExecProcess(parsed.cutSet, parsed.data, cmd, parsed.flagsLeft)

	p.argsLeft = parsed.argsLeft

	p.execList.Push(elem)

	return nil
}

// parseArgs parses the arguments and runs the command.
//
// Parameters:
//   - p: The program that the command is being executed on.
//   - args: The arguments that were passed to the command.
//   - cmd: The command to run.
//
// Returns:
//   - []string: The arguments that were left after parsing.
//   - error: An error if the command failed to execute.
func (p *Parser) parseArgs(cmd *cms.Command) error {
	if len(p.argsLeft) == 0 {
		err := p.handleCmd(cmd)
		if err != nil {
			return fmt.Errorf("in command %q: %w", cmd.Name, err)
		}

		return nil
	}

	subCmd := cmd.GetSubCommand(p.argsLeft[0])
	if subCmd != nil {
		// Execute the sub-command first and pass the result to the command.
		p.argsLeft = p.argsLeft[1:]

		err := p.parseArgs(subCmd)
		if err != nil {
			return fmt.Errorf("in sub-command %q: %w", subCmd.Name, err)
		}
	}

	// Handle the command.
	err := p.handleCmd(cmd)
	if err != nil {
		return fmt.Errorf("in command %q: %w", cmd.Name, err)
	}

	return nil
}
