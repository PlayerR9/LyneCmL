package Simple

import (
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

///////////////////////////////////////////////////////

// RunFunc is a function that will be executed when the command is called.
//
// Parameters:
//   - p: The program that the command is being executed on.
//   - data: The data that was passed to the command. (either the arguments
//     or the parsed arguments)
//
// Returns:
//   - error: An error if the command failed to execute.
type RunFunc func(p *Program, data any) error

var (
	// NoRunFunc is a function that does nothing.
	NoRunFunc RunFunc
)

func init() {
	NoRunFunc = func(p *Program, data any) error {
		return nil
	}
}

// DotMerging merges two slices of strings by using dot product.
//
// Parameters:
//   - lines: The first slice of strings.
//   - others: The second slice of strings.
//   - sep: The separator to use when merging the strings.
//
// Returns:
//   - []string: The merged strings.
//
// Example:
//   - lines: ["a", "b"]
//   - others: ["1", "2"]
//   - sep: "-"
//   - returns: ["a-1", "a-2", "b-1", "b-2"]
//
// Behaviors:
//   - If either lines or others is empty, the other slice will be returned.
func DotMerging(lines []string, others []string, sep string) []string {
	if len(others) == 0 {
		return lines
	} else if len(lines) == 0 {
		return others
	}

	var results []string

	var builder strings.Builder

	for _, line := range others {
		if line == "" {
			results = append(results, lines...)
			continue
		}

		for _, l := range lines {
			if l == "" {
				results = append(results, line)
			} else {
				builder.WriteString(l)
				builder.WriteString(sep)
				builder.WriteString(line)

				results = append(results, builder.String())
				builder.Reset()
			}
		}
	}

	return results
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
	Description []string

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
	flags map[string]*Flag
}

// GenerateUsage implements the CmlComponent interface.
func (c *Command) GenerateUsage() []string {
	lines := []string{c.Name}

	// Deal with sub-commands
	for _, cmd := range c.subCommands {
		usages := cmd.GenerateUsage()

		lines = DotMerging(lines, usages, " ")
	}

	// Deal with arguments
	argumentLines := c.Argument.GenerateUsage()
	lines = DotMerging(lines, argumentLines, " ")

	// Deal with flags
	var flagLines []string
	var builder strings.Builder

	for _, flag := range c.flags {
		str := strings.Join(flag.Usages, " | ")

		builder.WriteRune('[')
		builder.WriteString(str)
		builder.WriteRune(']')

		flagLines = append(flagLines, builder.String())
		builder.Reset()
	}

	lines = DotMerging(lines, flagLines, " ")

	return lines
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
func (c *Command) SetFlags(flags ...*Flag) {
	top := 0

	for i := 0; i < len(flags); i++ {
		if flags[i] != nil {
			flags[top] = flags[i]
			top++
		}
	}

	flags = flags[:top]
	if len(flags) == 0 {
		return
	}

	if c.flags == nil {
		c.flags = make(map[string]*Flag)
	}

	for _, flag := range flags {
		flagName := flag.LongName
		if flagName == "" {
			continue
		}

		_, ok := c.flags[flagName]
		if !ok {
			c.flags[flagName] = flag
		}
	}
}

type CommandFSSetting struct {
	Spacing string
}

func WithSpacing(spacing int) ffs.Option {
	if spacing <= 0 {
		spacing = 1
	}

	space := strings.Repeat(" ", spacing)

	return func(s ffs.Settinger) {
		setting, ok := s.(*CommandFSSetting)
		if !ok {
			return
		}

		setting.Spacing = space
	}
}

// FString implements the ffs.FStringer interface.
func (c *Command) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	settings := &CommandFSSetting{
		Spacing: " ",
	}

	for _, opt := range opts {
		opt(settings)
	}

	tab_size := trav.GetConfig().GetTabSize()

	// Command: <name>
	err := trav.AddJoinedLine(" ", "Command:", c.Name)
	if err != nil {
		return err
	}

	trav.EmptyLine()

	ta := fs.NewTableAligner()

	// Usage: <usage>
	if len(c.Usages) == 1 {
		err = trav.AddJoinedLine(" ", "Usage:", c.Usages[0])
		if err != nil {
			return err
		}
	} else {
		ta.SetHead("Usage:")

		for _, usage := range c.Usages {
			ta.AddRow(usage)
		}

		lines, _ := ta.Build(tab_size, true)

		err = trav.AddLines(lines)
		if err != nil {
			return err
		}
	}

	// Flags:
	// 	<flags>
	if len(c.flags) > 0 {
		ta.SetHead("Flags:")

		for _, flag := range c.flags {
			str := strings.Join(flag.Usages, settings.Spacing)

			ta.AddRow(str, flag.Brief)
		}

		lines, _ := ta.Build(tab_size, true)

		err = trav.AddLines(lines)
		if err != nil {
			return err
		}
	}

	if c.Description == nil {
		return nil
	}

	// Description:
	// 	<description>
	trav.EmptyLine()

	ta.SetHead("Description:")
	for _, line := range c.Description {
		ta.AddRow(line)
	}

	lines, _ := ta.Build(tab_size, true)

	err = trav.AddLines(lines)
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
		if cmd.Name == "" {
			continue
		}

		_, ok := c.subCommands[cmd.Name]
		if !ok {
			c.subCommands[cmd.Name] = cmd
		}
	}
}

// GetFlagMap gets the flags of the command.
//
// Returns:
//   - map[string]*Flag: The flags of the command.
func (c *Command) GetFlagMap() map[string]*Flag {
	return c.flags
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

	return flag.value, true
}

// GetSubCommand gets a sub-command of the command.
//
// Parameters:
//   - name: The name of the sub-command.
//
// Returns:
//   - *Command: The sub-command. Nil if not found.
func (c *Command) GetSubCommand(name string) *Command {
	if len(c.subCommands) == 0 {
		return nil
	}

	cmd, ok := c.subCommands[name]
	if !ok {
		return nil
	}

	return cmd
}
