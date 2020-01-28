package main

import (
	"flag"
	"fmt"
)

// Supported commands
const (
	add      = "add"
	child    = "child"
	create   = "create"
	dot      = "dot"
	generate = "generate"
	list     = "list"
	listall  = "listall"
	parent   = "parent"
	remove   = "rm"
	task     = "task"
)

type command struct {
	args   []string
	runner func() error
}

func newCommand(arg string) (*command, error) {
	c := &command{}
	if arg == task {
		c.args = flag.Args()[1:]
		c.runner = c.task
		
		return c, nil
	}

	if err := setTask(c, arg); err != nil {
		return nil, err
	}

	return c, nil
}

func setTask(c *command, arg string) error {
	switch arg {
	case list:
		c.runner = c.list
	case listall:
		c.runner = c.listall
	case dot:
		c.runner = c.dot
	case remove:
		if flag.NArg() < 2 {
			return fmt.Errorf(errNoArg, remove)
		}

		c.args = flag.Args()[1:]
		c.runner = c.rm
	case add:
		if flag.NArg() < 2 {
			return fmt.Errorf(errNoArg, add)
		}

		c.args = flag.Args()[1:]
		c.runner = c.add
	case generate:
		c.runner = c.generate
	default:
		return fmt.Errorf(errBadArg, arg)
	}

	return nil
}

func (c *command) add() error {
	if len(c.args) != 3 {
		return fmt.Errorf(errNoArg, add)
	}

	l, err := layoutFromTodoFile()
	if err != nil {
		return err
	}

	switch c.args[0] {
	case parent:
		l.addLink(c.args[2], c.args[1])
	case child:
		l.addLink(c.args[1], c.args[2])
	default:
		return fmt.Errorf(errBadArg, c.args[0])
	}

	writeTodo(l)

	return nil
}

func (c *command) dot() error {
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != create {
		return err
	}

	writeDot(l)

	return nil
}

// generate walks the file looking for a handful of known tokens
func (c *command) generate() error {
	g := &generator{}
	if g.dotTodoExists() {
		err := g.parseTodo()
		if err != nil {
			return err
		}
	}

	l := g.toLayout()
	writeTodo(l)

	return nil
}

func (c *command) list() error {
	l, err := layoutFromTodoFile()
	if err != nil {
		return err
	}

	writeList(l)

	return nil
}

func (c *command) listall() error {
	l, err := layoutFromTodoFile()
	if err != nil {
		return err
	}

	writeListAll(l)

	return nil
}

func (c *command) rm() error {
	if len(c.args) != 3 {
		return fmt.Errorf(errNoArg, remove)
	}

	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != create {
		return err
	}

	defer writeTodo(l)

	switch c.args[0] {
	case "parent":
		l.rmLink(c.args[2], c.args[1])
	case "child":
		l.rmLink(c.args[1], c.args[1])
	default:
		return fmt.Errorf(errBadArg, c.args[0])
	}

	return nil
}

func (c *command) run() error {
	if c.runner != nil {
		err := c.runner()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *command) task() error {
	if len(c.args) < 2 {
		return fmt.Errorf(errNoArg, task)
	}

	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != create {
		return err
	}

	switch c.args[0] {
	case create:
		err = l.create(c.args[1])
		if err != nil {
			return err
		}

		writeTodo(l)

		return nil
	case "destroy":
		err := l.destroy(c.args[1])
		if err != nil {
			return err
		}

		writeTodo(l)

		return nil
	case "add":
		if l.taskExists(c.args[1], c.args[2]) {
			return fmt.Errorf(errExists, c.args[1])
		}

		if e := l.addTask(c.args[1], c.args[2]); e != nil {
			return e
		}

		writeTodo(l)

		return nil
	case "rm":
		if !l.taskExists(c.args[1], c.args[2]) {
			return fmt.Errorf(errNoEntry, c.args[2])
		}

		defer writeTodo(l)

		return l.rmTask(c.args[1], c.args[2])
	case "toggle":
		if !l.taskExists(c.args[1], c.args[2]) {
			return fmt.Errorf(errNoEntry, c.args[1])
		}

		defer writeTodo(l)

		return l.toggleTask(c.args[1], c.args[2])
	default:
		return fmt.Errorf(errBadArg, c.args[0])
	}
}
