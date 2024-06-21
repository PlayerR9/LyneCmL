package Parser

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	util "github.com/PlayerR9/LyneCmL/Common/util"
	cms "github.com/PlayerR9/LyneCmL/Simple"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	utse "github.com/PlayerR9/MyGoLib/Utility/StringExt"
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
	flagList []any

	// args is the list of arguments to parse.
	args []string
}

// findLongFlag finds the index of the long flag.
//
// Parameters:
//   - long: The long flag to find.
//
// Returns:
//   - int: The index of the long flag. -1 if not found.
func (p *flagParser) findLongFlag(long string) int {
	searchLongFlagFunc := func(flag any) bool {
		fg := flag.(*cms.Flag[any])

		return fg.LongName == long
	}

	index := slices.IndexFunc(p.flagList, searchLongFlagFunc)

	return index
}

// parseSingleBoolLongFlag parses a single boolean long flag.
//
// Parameters:
//   - flag: The flag to parse.
//   - right: The right side of the flag.
//
// Returns:
//   - error: An error if the flag failed to parse.
func (p *flagParser) parseSingleBoolLongFlag(flag *cms.Flag[bool], right string) error {
	if !flag.HasArgument {
		if right != "" {
			return fmt.Errorf("flag --%s does not take an argument", flag.LongName)
		}

		flag.Value = true

		return nil
	}

	if right == "" {
		if len(p.args) == 0 {
			return fmt.Errorf("flag --%s requires an argument", flag.LongName)
		}

		right = p.args[0]

		p.args = p.args[1:]
	}

	value, err := flag.ParseFunc(right)
	if err != nil {
		return fmt.Errorf("in flag --%s: %w", flag.LongName, err)
	}

	flag.Value = value

	return nil
}

// parseSingleAnyLongFlag parses a single any long flag.
//
// Parameters:
//   - flag: The flag to parse.
//   - right: The right side of the flag.
//
// Returns:
//   - error: An error if the flag failed to parse.
func (p *flagParser) parseSingleAnyLongFlag(flag *cms.Flag[any], right string) error {
	if !flag.HasArgument {
		return fmt.Errorf("flag --%s is not a boolean flag", flag.LongName)
	}

	if right == "" {
		if len(p.args) == 0 {
			return fmt.Errorf("flag --%s requires an argument", flag.LongName)
		}

		right = p.args[0]
		p.args = p.args[1:]
	}

	value, err := flag.ParseFunc(right)
	if err != nil {
		return fmt.Errorf("in flag --%s: %w", flag.LongName, err)
	}

	flag.Value = value

	return nil
}

// parseSingleLongFlag parses a single long flag.
//
// Format:
//
//	--flag		: only boolean flags without arguments
//	--flag=value	: any flag that requires an argument
//	--flag value 	: any flag that requires an argument
//
// Parameters:
//   - header: The header of the flag.
//
// Returns:
//   - error: An error if the flag failed to parse.
func (p *flagParser) parseSingleLongFlag(header string) error {
	fields := strings.SplitN(header, "=", 2)

	left := fields[0]
	index := p.findLongFlag(left)
	if index == -1 {
		return fmt.Errorf("unknown flag --%s", left)
	}

	var right string

	if len(fields) == 2 {
		right = fields[1]
	} else {
		right = ""
	}

	switch flag := p.flagList[index].(type) {
	case *cms.Flag[bool]:
		err := p.parseSingleBoolLongFlag(flag, right)
		if err != nil {
			return err
		}
	case *cms.Flag[any]:
		err := p.parseSingleAnyLongFlag(flag, right)
		if err != nil {
			return err
		}
	}

	p.flagList = slices.Delete(p.flagList, index, index+1)

	return nil
}

// findShortFlag finds the index of the short flag.
//
// Parameters:
//   - short: The short flag to find.
//
// Returns:
//   - int: The index of the short flag. -1 if not found.
func (p *flagParser) findShortFlag(short rune) int {
	searchShortFlagFunc := func(flag any) bool {
		fg := flag.(*cms.Flag[any])

		return fg.ShortName == short
	}

	index := slices.IndexFunc(p.flagList, searchShortFlagFunc)

	return index
}

// parseSingleBoolShortFlag parses a single boolean short flag.
//
// Parameters:
//   - flag: The flag to parse.
//   - right: The right side of the flag.
//
// Returns:
//   - string: The right side of the flag.
//   - error: An error if the flag failed to parse.
func (p *flagParser) parseSingleBoolShortFlag(flag *cms.Flag[bool], right string) (string, error) {
	if !flag.HasArgument {
		flag.Value = true

		return right, nil
	}

	if right == "" {
		if len(p.args) == 0 {
			return right, fmt.Errorf("flag -%c requires an argument", flag.ShortName)
		}

		right = p.args[0]
		p.args = p.args[1:]
	}

	value, err := flag.ParseFunc(right)
	if err != nil {
		return "", fmt.Errorf("in flag -%c: %w", flag.ShortName, err)
	}

	flag.Value = value

	return "", nil
}

// parseSingleAnyShortFlag parses a single any short flag.
//
// Parameters:
//   - flag: The flag to parse.
//   - right: The right side of the flag.
//
// Returns:
//   - string: The right side of the flag.
//   - error: An error if the flag failed to parse.
func (p *flagParser) parseSingleAnyShortFlag(flag *cms.Flag[any], right string) (string, error) {
	if !flag.HasArgument {
		return right, fmt.Errorf("flag -%c is not a boolean flag", flag.ShortName)
	}

	if right == "" {
		if len(p.args) == 0 {
			return right, fmt.Errorf("flag -%c requires an argument", flag.ShortName)
		}

		right = p.args[0]
		p.args = p.args[1:]
	}

	value, err := flag.ParseFunc(right)
	if err != nil {
		return "", fmt.Errorf("in flag -%c: %w", flag.ShortName, err)
	}

	flag.Value = value

	return "", nil
}

// parseSingleShortFlag parses a single short flag.
//
// Format:
//
//	-flag		: only boolean flags without arguments
//	-flag=value	: any flag that requires an argument
//	-flag value 	: any flag that requires an argument
//	-abc		: multiple short flags
func (p *flagParser) parseSingleShortFlag(short rune, right string) (string, bool, error) {
	index := p.findShortFlag(short)
	if index == -1 {
		return right, false, nil
	}

	var err error

	switch flag := p.flagList[index].(type) {
	case *cms.Flag[bool]:
		right, err = p.parseSingleBoolShortFlag(flag, right)
	case *cms.Flag[any]:
		right, err = p.parseSingleAnyShortFlag(flag, right)
	}
	if err != nil {
		return right, false, err
	}

	p.flagList = slices.Delete(p.flagList, index, index+1)

	return right, true, nil
}

// parseShortFlags parses the short flags.
//
// Parameters:
//   - header: The header of the flags.
//
// Returns:
//   - error: An error if the flags failed to parse.
func (p *flagParser) parseShortFlags(header string) error {
	fields := strings.SplitN(header, "=", 2)

	header = fields[0]

	var right string

	if len(fields) == 2 {
		right = fields[1]
	} else {
		right = ""
	}

	runes, err := utse.ToUTF8Runes(header)
	if err != nil {
		return fmt.Errorf("error converting to runes: %w", err)
	}

	if len(runes) == 1 {
		var ok bool

		right, ok, err = p.parseSingleShortFlag(runes[0], right)
		if err != nil {
			return err
		} else if !ok {
			return nil
		}

		if right != "" {
			return fmt.Errorf("extra argument %q", right)
		}
	} else {
		var ok bool

		for _, letter := range runes {
			right, ok, err = p.parseSingleShortFlag(letter, right)
			if err != nil {
				return err
			} else if !ok {
				return fmt.Errorf("invalid merged short flags -%s", header)
			}
		}

		if right != "" {
			return fmt.Errorf("extra argument %q", right)
		}
	}

	return nil
}

// parseOne parses one flag.
//
// Returns:
//   - bool: True if the flag was parsed. False otherwise.
//   - error: An error if the flag failed to parse.
func (fp *flagParser) parseOne() (bool, error) {
	header := fp.args[0]
	fp.args = fp.args[1:]

	ok := strings.HasPrefix(header, LongFlagPrefix)
	if ok {
		header = strings.TrimPrefix(header, LongFlagPrefix)

		err := fp.parseSingleLongFlag(header)
		if err != nil {
			return false, fmt.Errorf("in long flag --%s: %w", header, err)
		}

		return true, nil
	}

	ok = strings.HasPrefix(header, ShortFlagPrefix)
	if !ok {
		return false, nil
	}

	header = strings.TrimPrefix(header, ShortFlagPrefix)

	err := fp.parseShortFlags(header)
	if err != nil {
		return false, fmt.Errorf("in short flag -%s: %w", header, err)
	}

	return true, nil
}

type parsedData struct {
	cutSet []string
	data   any
	i      int
}

// validate is a helper function that validates the number of arguments.
//
// Parameters:
//   - args: The arguments to validate.
//
// Returns:
//   - []string: The arguments if they are valid.
//   - error: An error if the arguments are invalid.
func validate(a *cms.Argument, args []string) (*parsedData, error) {
	left := a.GetMin()

	if len(args) < left {
		return nil, util.NewErrFewArguments(left, len(args))
	}

	right := a.GetMax()

	if right == -1 {
		right = len(args)
	}

	var data any
	var err error

	for i := right; i >= left; i-- {
		cutSet := args[:i]

		data, err = a.Apply(cutSet)
		if err == nil {
			return &parsedData{
				cutSet: cutSet,
				data:   data,
				i:      i,
			}, nil
		}
	}

	return nil, fmt.Errorf("error parsing arguments: %w", err)
}

type parsingResult struct {
	cutSet   []string
	data     any
	argsLeft []string
}

func parseFlags(flagMap map[string]any, args []string) ([]string, error) {
	var flagList []any

	for _, value := range flagMap {
		flagList = append(flagList, value)
	}

	fp := &flagParser{
		flagList: flagList,
		args:     args,
	}

	for len(fp.args) > 0 {
		ok, err := fp.parseOne()
		if err != nil {
			return nil, err
		} else if !ok {
			return fp.args, nil
		}

		if len(fp.flagList) == 0 {
			return fp.args, nil
		}
	}

	for _, flag := range fp.flagList {
		switch flag := flag.(type) {
		case *cms.Flag[bool]:
			if flag.HasArgument {
				flag.Value = flag.DefaultVal
			} else {
				flag.Value = false
			}
		case *cms.Flag[any]:
			flag.Value = flag.DefaultVal
		}
	}

	return nil, nil
}

func parseArguments(a *cms.Argument, flagMap map[string]any, args []string) (*parsingResult, error) {
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

	parsed, err := validate(a, actualArgs)
	if err != nil {
		return nil, err
	}

	remaining, err := parseFlags(flagMap, args[parsed.i:])
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	result := &parsingResult{
		cutSet:   parsed.cutSet,
		data:     parsed.data,
		argsLeft: remaining,
	}

	return result, nil
}

func ParseArgs(commandMap map[string]*cms.Command, args []string) (lls.Stacker[*cms.ExecProcess], error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments")
	}

	p := &Parser{
		argsLeft: args,
		execList: lls.NewArrayStack[*cms.ExecProcess](),
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

	execList *lls.ArrayStack[*cms.ExecProcess]
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

	elem := cms.NewExecProcess(parsed.cutSet, parsed.data, cmd)

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
