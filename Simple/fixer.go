package Simple

import (
	"fmt"
	"strings"
	"unicode"

	cut "github.com/PlayerR9/LyneCmL/Common/util"
	pd "github.com/PlayerR9/LyneCmL/Simple/display"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	utse "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

// isValidLetterName checks if the character is a valid letter name.
//
// Parameters:
//   - char: The character to check.
//
// Returns:
//   - bool: True if the character is a valid letter name. False otherwise.
func isValidLetterName(char rune) bool {
	switch char {
	case '-', '_':
		return true
	default:
		ok := unicode.IsLetter(char)
		if !ok {
			return false
		}
	}

	return true
}

// fixName fixes the name by removing invalid characters.
//
// Parameters:
//   - str: The string to fix.
//
// Returns:
//   - string: The fixed string.
//   - error: An error if the string failed to fix.
func fixName(str string) (string, error) {
	runes, err := utse.ToUTF8Runes(str)
	if err != nil {
		return "", err
	}

	left := -1

	for i := 0; i < len(runes); i++ {
		ok := isValidLetterName(runes[i])
		if ok && runes[i] != '-' {
			left = i
			break
		}
	}

	if left == -1 {
		return "", nil
	}

	right := -1

	for i := len(runes) - 1; i >= 0; i-- {
		ok := isValidLetterName(runes[i])
		if ok && runes[i] != '-' {
			right = i + 1
			break
		}
	}

	if right == -1 {
		return "", nil
	}

	runes = runes[:right]

	for i, char := range runes {
		ok := isValidLetterName(char)
		if !ok {
			return "", fmt.Errorf("invalid character at index %d: %c", i+left, char)
		}
	}

	return string(runes), nil

}

// fixSpacing fixes the spacing by removing extra spaces.
//
// Parameters:
//   - str: The string to fix.
//
// Returns:
//   - string: The fixed string.
func fixSpacing(str string) string {
	fields := strings.Fields(str)
	return strings.Join(fields, " ")
}

// fixFlags fixes the flags by removing invalid flags and checking for conflicts.
//
// Parameters:
//   - flagMap: The flags to fix.
//
// Returns:
//   - map[string]*Flag: The fixed flags.
//   - error: An error if the flags failed to fix.
func fixFlags(flagMap map[string]*Flag) (map[string]*Flag, error) {
	fixedFlags := make(map[string]*Flag)

	if len(flagMap) == 0 {
		return fixedFlags, nil
	}

	for name, flag := range flagMap {
		err := flag.Fix()
		if err != nil {
			return nil, fmt.Errorf("invalid flag %s: %w", name, err)
		}

		if flag.LongName != "" {
			fixedFlags[name] = flag
		}
	}

	// Check conflicts between the flags.
	seen := make(map[string]*Flag)

	for _, flag := range fixedFlags {
		var tocheck []string

		if flag.ShortName != 0 {
			tocheck = append(tocheck, string(flag.ShortName))
		}

		tocheck = append(tocheck, flag.LongName)

		for _, name := range tocheck {
			prev, ok := seen[name]
			if ok {
				var reason error

				if prev.LongName == flag.LongName {
					reason = fmt.Errorf("duplicate flag long name")
				} else {
					reason = fmt.Errorf("duplicate flag short name")
				}

				return nil, cut.NewErrFlagConflict(flag.LongName, prev.LongName, reason)
			}

			seen[name] = flag
		}
	}

	return fixedFlags, nil
}

// Fix implements the CmlComponent interface.
//
// This never errors.
func (a *Argument) Fix() error {
	// bounds do not need to be fixed.

	// Fix: parse function
	if a.parseFunc == nil {
		a.parseFunc = DefaultParseFunc
	}

	return nil
}

// Fix implements Flager interface.
func (f *Flag) Fix() error {
	// Fix: long name
	longName, err := fixName(f.LongName)
	if err != nil {
		return fmt.Errorf("while fixing long name: %w", err)
	}

	f.LongName = longName

	// Fix: short name
	if f.ShortName != 0 {
		ok := isValidLetterName(f.ShortName)
		if !ok {
			return fmt.Errorf("invalid short name: %c", f.ShortName)
		}

		if f.ShortName == '_' {
			f.ShortName = 0
		}
	}

	// Fix: brief
	f.Brief = fixSpacing(f.Brief)

	// Fix: usages
	if f.Usages != nil {
		for i := 0; i < len(f.Usages); i++ {
			fixed := fixSpacing(f.Usages[i])
			f.Usages[i] = fixed
		}

		f.Usages = us.RemoveEmpty(f.Usages)
	}

	if len(f.Usages) == 0 {
		newUsage := f.GenerateUsage()
		f.Usages = newUsage
	}

	// Fix: description
	if f.Description != nil {
		isEmpty := true

		for i := 0; i < len(f.Description); i++ {
			if f.Description[i] != "" {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			f.Description = nil
		}
	}

	// value does not need to be fixed.

	// argument does not need to be fixed.

	return nil
}

// Fix implements the CmlComponent interface.
func (c *Command) Fix() error {
	// Fix: name
	name, err := fixName(c.Name)
	if err != nil {
		return fmt.Errorf("while fixing name: %w", err)
	}

	c.Name = name

	// Fix: usages
	if c.Usages != nil {
		for i := 0; i < len(c.Usages); i++ {
			fixed := fixSpacing(c.Usages[i])
			c.Usages[i] = fixed
		}

		c.Usages = us.RemoveEmpty(c.Usages)
	}

	if len(c.Usages) == 0 {
		newUsage := c.GenerateUsage()
		c.Usages = newUsage
	}

	// Fix: brief
	c.Brief = fixSpacing(c.Brief)

	// Fix: description
	if c.Description != nil {
		isEmpty := true

		for i := 0; i < len(c.Description); i++ {
			if c.Description[i] != "" {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			c.Description = nil
		}
	}

	// Fix: argument
	if c.Argument == nil {
		c.Argument = NoArgument
	}

	c.Argument.Fix()

	// Fix: run function
	if c.Run == nil {
		c.Run = NoRunFunc
	}

	// Fix: subCommands
	if c.subCommands != nil {
		fixedSubCommands := make(map[string]*Command)

		for name, subCmd := range c.subCommands {
			err := subCmd.Fix()
			if err != nil {
				return fmt.Errorf("invalid sub-command %s: %w", subCmd.Name, err)
			}

			if subCmd.Name != "" {
				fixedSubCommands[name] = subCmd
			}
		}

		c.subCommands = fixedSubCommands
	}

	// Fix: flags
	flagMap, err := fixFlags(c.flags)
	if err != nil {
		return fmt.Errorf("while fixing flags: %w", err)
	}

	c.flags = flagMap

	return nil
}

// Fix implements the CmlComponent interface.
//
// This never errors.
func (p *Program) Fix() error {
	// Fix: name
	name, err := fixName(p.Name)
	if err != nil {
		return fmt.Errorf("while fixing name: %w", err)
	}

	p.Name = name

	// Fix: brief
	p.Brief = fixSpacing(p.Brief)

	// Fix: description
	if p.Description != nil {
		hasEmpty := true

		for i := 0; i < len(p.Description); i++ {
			if p.Description[i] != "" {
				hasEmpty = false
				break
			}
		}

		if hasEmpty {
			p.Description = nil
		}
	}

	// Fix: version
	p.Version = fixSpacing(p.Version)

	// Fix: commands
	if p.commands == nil {
		p.commands = make(map[string]*Command)
	}

	if len(p.commands) != 0 {
		fixedCmds := make(map[string]*Command)

		for _, cmd := range p.commands {
			err := cmd.Fix()
			if err != nil {
				return fmt.Errorf("invalid command %s: %w", cmd.Name, err)
			}

			if cmd.Name != "" {
				fixedCmds[cmd.Name] = cmd
			}
		}

		p.commands = fixedCmds
	} else {
		p.commands = make(map[string]*Command)
	}

	p.commands[cut.HelpCmdOpcode] = HelpCmd

	// Display does not need to be fixed.

	// Fix: flags
	flagMap, err := fixFlags(p.flags)
	if err != nil {
		return fmt.Errorf("while fixing flags: %w", err)
	}

	p.flags = flagMap

	// Fix: configTable
	if p.Configs == nil {
		p.Configs = pd.NewDisplayConfigs()
	}

	err = p.LoadConfigs()
	if err != nil {
		return fmt.Errorf("while loading configs: %w", err)
	}

	return nil
}
