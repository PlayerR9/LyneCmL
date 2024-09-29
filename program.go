package cml

import "fmt"

type Program struct {
	commands map[string]*Command
}

func (p *Program) AddCommand(cmd *Command) {
	if p == nil || cmd == nil {
		return
	}

	p.commands[cmd.Name] = cmd
}

func (p *Program) AddCommands(cmds ...*Command) {
	if p == nil {
		return
	}

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		p.commands[cmd.Name] = cmd
	}
}

func (p Program) Run(args []string) error {
	if len(args) < 2 {
		fmt.Println("Use \"help\" for a list of commands")

		return fmt.Errorf("no command specified")
	}

	command := args[1]

	cmd, ok := p.commands[command]
	if !ok {
		fmt.Println("Use \"help\" for a list of commands")

		return fmt.Errorf("unknown command: %s", command)
	}

	err := cmd.RunFn()
	return err
}

func DefaultExitSequence(err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}

	fmt.Println()

	fmt.Println("Press ENTER to exit...")

	fmt.Scanln()
}
