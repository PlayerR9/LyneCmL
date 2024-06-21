package Simple

import (
	"strings"

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

	// Usages are the usages of the command.
	Usages []string

	// Brief is a brief description of the command.
	Brief string

	// Description is a description of the command.
	Description *DescBuilder

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

	// flags is a map of flags that the command can execute.
	flags map[string]any
}

// GenerateUsage implements the CmlComponent interface.
func (c *Command) GenerateUsage() []string {
	var lines []string

	// Deal with arguments

	argumentLines := c.Argument.GenerateUsage()

	var builder strings.Builder

	for _, line := range argumentLines {
		if line == "" {
			lines = append(lines, c.Name)
		} else {
			builder.WriteString(c.Name)
			builder.WriteRune(' ')
			builder.WriteString(line)

			lines = append(lines, builder.String())
			builder.Reset()
		}
	}

	/*
		// name args [flags]

		// Name is the name of the command.
		Name string


		// Argument is the argument of the command.
		Argument *Argument


		// subCommands is a map of sub-commands that the command can execute.
		// If at least one sub-command is present, then the first argument
		// of the command will be the sub-command.
		//
		// Then the return value of the run function will be passed to the
		// sub-command as the data argument.
		subCommands map[string]*Command

		// flags is a map of flags that the command can execute.
		flags map[string]Flager[any]
	*/

	return lines
}

// Fix implements the CmlComponent interface.
func (c *Command) Fix() {
	// Fix argument
	if c.Argument == nil {
		c.Argument = NoArgument
	}

	c.Argument.Fix()

	// Fix command
	c.Name = strings.TrimSpace(c.Name)

	if c.Usages != nil {
		for i := 0; i < len(c.Usages); i++ {
			c.Usages[i] = strings.TrimSpace(c.Usages[i])
		}

		c.Usages = us.RemoveEmpty(c.Usages)
	}

	if len(c.Usages) == 0 {
		newUsage := c.GenerateUsage()
		c.Usages = newUsage
	}

	c.Brief = strings.TrimSpace(c.Brief)

	if c.Run == nil {
		c.Run = NoRunFunc
	}
}

// SetFlags sets the flags of the command.
//
// Parameters:
//   - flags: The flags to set.
//
// Behaviors:
//   - If a flag has the same name as another flag, the first one will be kept.
//   - nil or empty flags will be filtered out.
//   - Flags are fixed before being set.
func (c *Command) SetFlags(flags ...any) {
	if c.flags == nil {
		c.flags = make(map[string]any)
	}

	for _, flag := range flags {
		if flag == nil {
			continue
		}

		switch flag := flag.(type) {
		case *Flag[any]:
			flag.Fix()

			flagName := flag.GetName()
			if flagName == "" {
				continue
			}

			_, ok := c.flags[flagName]
			if !ok {
				c.flags[flagName] = flag
			}
		case *BoolFlag:
			flag.Fix()

			flagName := flag.GetName()
			if flagName == "" {
				continue
			}

			_, ok := c.flags[flagName]
			if !ok {
				c.flags[flagName] = flag
			}
		}
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
	if len(c.Usages) == 1 {
		err = trav.AddJoinedLine(" ", "Usage:", c.Usages[0])
		if err != nil {
			return err
		}
	} else {
		err = trav.AddLine("Usage:")
		if err != nil {
			return err
		}

		for _, usage := range c.Usages {
			err := ffs.ApplyFormFunc(
				trav.GetConfig(
					ffs.WithModifiedIndent(1),
				),
				trav,
				usage,
				func(trav *ffs.Traversor, elem string) error {
					err := trav.AddLine(elem)
					if err != nil {
						return err
					}

					return nil
				},
			)
			if err != nil {
				return err
			}
		}
	}

	// Flags:
	// 	<flags>
	if c.flags != nil {
		glTableAligner.SetHead("Flags:")

		for name, flag := range c.flags {
			var usage string

			switch flag := flag.(type) {
			case *BoolFlag:
				usage = flag.Usage
			case *Flag[any]:
				usage = flag.Usage
			}

			glTableAligner.AddRow([]string{name, usage})
		}

		err = glTableAligner.FString(trav)
		if err != nil {
			return err
		}

		glTableAligner.Reset()
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
//   - Sub-commands are fixed before being added.
func (c *Command) AddSubCommand(cmds ...*Command) {
	cmds = us.FilterNilValues(cmds)

	if len(cmds) == 0 {
		return
	}

	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}

	for _, cmd := range cmds {
		cmd.Fix()
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

// GetFlag gets the value of a flag.
//
// Parameters:
//   - name: The name of the flag.
//
// Returns:
//   - any: The value of the flag.
//   - bool: true if the flag was found, false otherwise.
func (c *Command) GetFlag(name string) (any, bool) {
	if c.flags == nil {
		return nil, false
	}

	flag, ok := c.flags[name]
	if !ok {
		return nil, false
	}

	switch flag := flag.(type) {
	case *Flag[any]:
		value := flag.GetValue()

		return value, true
	case *BoolFlag:
		value := flag.GetValue()

		return value, true
	}

	panic("unreachable")
}
