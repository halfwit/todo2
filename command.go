package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/altid/fslib"
)

// Initialize a command here to run in main
// Generally we could just run a command based on a case match
// but should this ever need to scale, this design will help us
type command struct {
	mkfile string
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
	if err := c.setEnv(); err != nil {
		return nil, err
	}
	if err := c.setTask(arg); err != nil {
		return nil, err
	}
	return c, nil
}

// This is pure bloat at the moment, a hashmap would be fine for mapping args to funcs
// as would calling a func itself in the case match. We want to future proof against
// needing to do any additional bookkeeping in the functions themselves
func (c command) setTask(arg string) error {
	switch arg {
	//case "help":
		//c.runner = help
	case "init":
		c.runner = initFile
	case "list":
		c.runner = list
	case "listall":
		c.runner = listall
	case "dot":
		c.runner = dot
	case "rm":
		if flag.NArg < 2 {
			return errors.New("No arguments supplied to rm")
		}
		c.args = flag.Args()[1:]
		c.runner = rm
	case "add":
		if flag.NArg < 2 {
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

func (c command) setEnv() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	data, err := fslib.UserShareDir()
	if err != nil {
		return err
	}
	c.mkfile = path.Join(data, "todo", path.Base(dir))
	return nil
}

// Broken out from main because there may be more busywork as things progress
func (c *command) runCmd() {
	if c.runner != nil {
		c.runner(c)
	}
}
