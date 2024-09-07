package simple

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// CmdRunFunc is a function that runs a command.
type CmdRunFunc func(p *Program, args []string) error

// Command is a struct that represents a command.
type Command struct {
	// Name is the name of the command.
	Name string

	// Brief is the brief of the command. Leave empty if not needed.
	Brief string

	// Description is the description of the command. Leave empty if not needed.
	Description []string

	// RunFunc is the function that runs the command.
	RunFunc CmdRunFunc

	// Argument is the argument of the command. If nil, NoArguments will be used.
	Argument *Argument

	// flag_list is the list of flags of the command.
	flag_list []*Flag
}

// HelpLines is a method that returns the help lines of the command.
//
// Returns:
//   - []string: The help lines of the command.
func (c Command) HelpLines() []string {
	var lines []string

	if c.Brief != "" {
		lines = append(lines, c.Name+" — "+c.Brief)
	} else {
		lines = append(lines, c.Name)
	}

	lines = append(lines, "")

	if len(c.Description) > 0 {
		lines = append(lines, c.Description...)
		lines = append(lines, "", "")
	}

	lines = append(lines, "Usage:")

	arr := c.Usage()

	if arr[1] != "" {
		lines = append(lines, arr[0]+" "+arr[1])
	} else {
		lines = append(lines, arr[0])
	}

	return nil
}

// Fix is a method that fixes the command.
func (c *Command) Fix() error {
	if c.Name == "" {
		return errors.New("command name cannot be empty")
	}

	if c.RunFunc == nil {
		c.RunFunc = func(_ *Program, _ []string) error {
			return nil
		}
	}

	if c.Argument == nil {
		c.Argument = NoArguments
	}

	return nil
}

// Usage is a method that returns the usage of the command.
//
// Returns:
//   - [2]string: The usage of the command and the brief of the command.
func (c Command) Usage() [2]string {
	arg := c.Argument.String()

	if arg == "" {
		return [2]string{c.Name, c.Brief}
	}

	return [2]string{c.Name + " " + arg, c.Brief}
}

func (c *Command) AddFlag(flag *Flag) error {
	if flag == nil {
		return nil
	}

	err := flag.Fix()
	if err != nil {
		return err
	}

	if flag.ShortName != 0 {
		for i := 0; i < len(c.flag_list); i++ {
			if c.flag_list[i].ShortName == flag.ShortName {
				return errors.New("flag already exists")
			}
		}
	}

	if flag.LongName != "" {
		for i := 0; i < len(c.flag_list); i++ {
			if c.flag_list[i].LongName == flag.LongName {
				return errors.New("flag already exists")
			}
		}
	}

	c.flag_list = append(c.flag_list, flag)

	return nil
}

// parse is a method that parses the command.
//
// Parameters:
//   - args: The arguments of the command.
//
// Returns:
//   - []string: The parsed arguments.
//   - error: An error if the arguments are invalid.
func (c Command) parse(args []string) ([]string, error) {
	parsed, err := c.Argument.check(args)
	if err != nil {
		return nil, err
	}

	left_args := args[len(parsed):]

	for _, arg := range left_args {
		for i := 0; i < len(c.flag_list); i++ {
			flag := c.flag_list[i]

			if strings.HasPrefix(arg, "--") {
				name := strings.TrimPrefix(arg, "--")

				if name == flag.LongName {
					flag.value = true

					c.flag_list[i] = flag

					continue
				}
			} else if strings.HasPrefix(arg, "-") {
				name := strings.TrimPrefix(arg, "-")

				if name == "" {
					continue
				}

				r, _ := utf8.DecodeRuneInString(name)
				if r == utf8.RuneError {
					return nil, errors.New("invalid flag")
				}

				if r == flag.ShortName {
					flag.value = true

					c.flag_list[i] = flag

					continue
				}
			} else {
				continue
			}
		}
	}

	return parsed, nil
}

func (c Command) ValueOf(name string) (bool, error) {
	for i := 0; i < len(c.flag_list); i++ {
		flag := c.flag_list[i]

		if flag.Name == name {
			return flag.value, nil
		}
	}

	return false, errors.New("flag not found")
}
