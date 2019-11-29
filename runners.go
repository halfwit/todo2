package main

import (
	"errors"
	"fmt"
)

func list(c *command) error {
	l, err := layoutFromTodoFile()
	if err != nil {
		return err
	}
	//d := dagFromLayout(l)
	//writeList(d)
	writeList(l)
	return nil
}

func listall(c *command) error {
	l, err := layoutFromTodoFile()
	if err != nil {
		return err
	}
	//d := dagFromLayout(l)
	//writeListAll(d)
	writeListAll(l)
	return nil
}

func dot(c *command) error {
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != "create" {
		return err
	}
	writeDot(l)
	return nil
}

func rm(c *command) error {
	if len(c.args) != 3 {
		return fmt.Errorf("Incorrect arguments supplied to rm")
	}
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != "create" {
		return err
	}
	defer writeTodo(l)
	switch c.args[0] {
	case "parent":
		l.rmLink(c.args[2], c.args[1])
	case "child":
		l.rmLink(c.args[1], c.args[1])
	default:
		return fmt.Errorf("Command not supported: %v", c.args[0])
	}
	return nil
}

func add(c *command) error {
	if len(c.args) != 3 {
		return fmt.Errorf("Incorrect arguments supplied to add")
	}
	l, err := layoutFromTodoFile()
	if err != nil {
		return err
	}
	switch c.args[0] {
	case "parent":
		l.addLink(c.args[2], c.args[1])
	case "child":
		l.addLink(c.args[1], c.args[2])
	default:
		return fmt.Errorf("Command not supported: %v %v", c.args[0], c.args[1])
	}
	writeTodo(l)
	return nil
}

// generate walks the file looking for a handful of known tokens
func generate(c *command) error {
	g := newGenerator()
	if g.dotTodoExists() {
		g.parseTodo()
	}
	l := g.toLayout()
	writeTodo(l)
	return nil
}

func task(c *command) error {
	if len(c.args) < 2 {
		return errors.New("Too few arguments supplied")
	}
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != "create" {
		return err
	}
	switch c.args[0] {
	case "create":
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
			return fmt.Errorf("Entry exists")
		}
		l.addTask(c.args[1], c.args[2])
		writeTodo(l)
		return nil
	case "rm":
		if !l.taskExists(c.args[1], c.args[2]) {
			return fmt.Errorf("No entry found")
		}
		defer writeTodo(l)
		return l.rmTask(c.args[1], c.args[2])
	case "toggle":
		if !l.taskExists(c.args[1], c.args[2]) {
			return fmt.Errorf("No such task/entry")
		}
		defer writeTodo(l)
		return l.toggleTask(c.args[1], c.args[2])
	default:
		return fmt.Errorf("Command not supported: %v", c.args[0])
	}
}
