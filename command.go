package main

import (
	"errors"
	"flag"
	"fmt"
)

// Supported commands
const (
	add      = "add"
	child    = "child"
	create   = "create"
	destroy  = "destroy"
	dot      = "dot"
	generate = "generate"
	list     = "list"
	listall  = "listall"
	parent   = "parent"
	remove   = "rm"
	task     = "task"
	toggle   = "toggle"
)

type command struct {
	args   []string
	runner func() (string, error)
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

func (c *command) add() (string, error) {
	if len(c.args) != 3 {
		return "", fmt.Errorf(errNoArg, add)
	}

	l, err := layoutFromTodoFile()
	if err != nil {
		return "", err
	}

	switch c.args[0] {
	case parent:
		l.addLink(c.args[2], c.args[1])
	case child:
		l.addLink(c.args[1], c.args[2])
	default:
		return "", fmt.Errorf(errBadArg, c.args[0])
	}

	return writeTodo(l)
}

func (c *command) dot() (string, error) {
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != create {
		return "", err
	}

	return writeDot(l)
}

// generate walks the file looking for a handful of known tokens
func (c *command) generate() (string, error) {
	g := &generator{}
	if g.dotTodoExists() {
		err := g.parseTodo()
		if err != nil {
			return "", err
		}
	}

	l := g.toLayout()

	return writeTodo(l)
}

func (c *command) list() (string, error) {
	l, err := layoutFromTodoFile()
	if err != nil {
		return "", err
	}

	return writeList(l)
}

func (c *command) listall() (string, error) {
	l, err := layoutFromTodoFile()
	if err != nil {
		return "", err
	}

	return writeListAll(l)
}

func (c *command) rm() (string, error) {
	if len(c.args) != 3 {
		return "", fmt.Errorf(errNoArg, remove)
	}

	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != create {
		return "", err
	}

	switch c.args[0] {
	case parent:
		l.rmLink(c.args[2], c.args[1])
	case child:
		l.rmLink(c.args[1], c.args[1])
	default:
		return "", fmt.Errorf(errBadArg, c.args[0])
	}

	return writeTodo(l)
}

func (c *command) run() (string, error) {
	if c.runner != nil {
		return c.runner()
	}

	return "", errors.New(errNoCmd)
}

func (c *command) task() (string, error) {
	if len(c.args) < 2 {
		return "", fmt.Errorf(errNoArg, task)
	}

	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != create {
		return "", err
	}

	switch c.args[0] {
	case create:
		err = l.create(c.args[1])
		if err != nil {
			return "", err
		}

		writeTodo(l)

		return "", nil
	case destroy:
		err := l.destroy(c.args[1])
		if err != nil {
			return "", err
		}

		writeTodo(l)

		return "", nil
	case add:
		if l.taskExists(c.args[1], c.args[2]) {
			return "", fmt.Errorf(errExists, c.args[1])
		}

		if e := l.addTask(c.args[1], c.args[2]); e != nil {
			return "", e
		}

		writeTodo(l)

		return "", nil
	case remove:
		if !l.taskExists(c.args[1], c.args[2]) {
			return "", fmt.Errorf(errNoEntry, c.args[2])
		}

		if e := l.rmTask(c.args[1], c.args[2]); e != nil {
			return "", e
		}

		return writeTodo(l)
	case toggle:
		if !l.taskExists(c.args[1], c.args[2]) {
			return "", fmt.Errorf(errNoEntry, c.args[1])
		}

		if e := l.toggleTask(c.args[1], c.args[2]); e != nil {
			return "", e
		}

		return writeTodo(l)
	default:
		return "", fmt.Errorf(errBadArg, c.args[0])
	}
}
