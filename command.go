package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/altid/fslib"
)

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

func (c command) setEnv() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	data, err := fslib.UserShareDir()
	if err != nil {
		return err
	}
	os.MkdirAll(path.Join(data, "todo"), 0755)
	c.mkfile = path.Join(data, "todo", path.Base(dir))
	return nil
}

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
		if flag.NArg() < 2 {
			return errors.New("No arguments supplied to rm")
		}
		c.args = flag.Args()[:1]
		c.runner = rm
	case "add":
		if flag.NArg() < 2 {
			return errors.New("No arguments supplied to add")
		}
		c.args = flag.Args()[:1]
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
