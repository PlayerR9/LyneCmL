package simple

import (
	"fmt"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
)

type CmdRunFn func(p *Program, args []string) error

type Command struct {
	Name     string
	Brief    string
	RunFn    CmdRunFn
	Argument *Argument
}

func (c *Command) Fix() error {
	if c == nil {
		return gcers.NilReceiver
	}

	name := strings.TrimSpace(c.Name)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	c.Name = name

	c.Brief = strings.TrimSpace(c.Brief)

	if c.RunFn == nil {
		c.RunFn = func(_ *Program, _ []string) error {
			return nil
		}
	}

	if c.Argument == nil {
		c.Argument = NoArguments
	} else {
		err := gcers.Fix("argument", c.Argument, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Command) parse(args []string) ([]string, error) {
	return c.Argument.parse(args)
}
