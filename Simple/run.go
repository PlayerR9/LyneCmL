package Simple

import (
	"errors"
	"fmt"
	"sync"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
)

func parseArgs(p *Program, args []string) (*Command, []string, error) {
	if len(args) < 1 {
		return nil, nil, errors.New("no command provided")
	}

	command, ok := p.commands[args[0]]
	if !ok {
		return nil, nil, fmt.Errorf("unknown command: %s", args[0])
	}

	if command.Argument == nil {
		command.Argument = NoArgument
	}

	parsedArgs, err := command.Argument.validate(args[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}

	return command, parsedArgs, nil
}

func runBody(p *Program, args []string, opts ...RunModeOption) error {
	command, parsedArgs, err := parseArgs(p, args[1:])
	if err != nil {
		return fmt.Errorf("could not parse arguments: %w", err)
	}

	p.buffer = sfb.NewBuffer[string]()

	p.buffer.Start()
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			str, ok := p.buffer.Receive()
			if !ok {
				break
			}

			fmt.Print(str)
		}
	}()

	err = command.Run(p, parsedArgs)

	p.buffer.Close()
	wg.Wait()

	if err != nil {
		return fmt.Errorf("could not run command: %w", err)
	}

	return nil
}

func (p *Program) Run(args []string, opts ...RunModeOption) {
	p.AddCommands(HelpCommandCmd)

	err := runBody(p, args, opts...)
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Program finished successfully")
	}
	fmt.Println()

	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()
}
