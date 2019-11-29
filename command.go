package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

type command struct {
	args   []string
	runner func(c *command) error
}

func newCommand(arg string) (*command, error) {
	c := &command{}

	if arg == "task" {
		c.args = flag.Args()[1:]
		c.runner = task
		return c, nil
	}
	if err := c.setTask(arg); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *command) setTask(arg string) error {
	switch arg {
	case "list":
		c.runner = list
	case "listall":
		c.runner = listall
	case "dot":
		c.runner = dot
	case "rm":
		if flag.NArg() < 2 {
			return errors.New("No arguments supplied to rm")
		}
		c.args = flag.Args()[1:]
		c.runner = rm
	case "add":
		if flag.NArg() < 2 {
			return errors.New("No arguments supplied to add")
		}
		c.args = flag.Args()[1:]
		c.runner = add
	case "generate":
		c.runner = generate
	default:
		return fmt.Errorf("Unknown command %q", arg)
	}
	return nil
}

func (c *command) run() {
	if c.runner != nil {
		err := c.runner(c)
		if err != nil {
			log.Fatal(err)
		}
	}
}
