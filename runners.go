package main

import (
	"errors"
	"fmt"
)

func list(c *command) error {
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != "create" {
		return err
	}
	writeList(l)
	return nil
}

func listall(c *command) error {
	l, err := layoutFromTodoFile()
	if err != nil && c.args[0] != "create" {
		return err
	}
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
		l.rmLink(parent, c.args[1], c.args[2])
	case "child":
		l.rmLink(child, c.args[1], c.args[2])
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
	if err != nil && c.args[0] != "create" {
		return err
	}
	switch c.args[0] {
	case "parent":
		l.addLink(parent, c.args[1], c.args[2])
	case "child":
		l.addLink(child, c.args[1], c.args[2])
	default:
		return fmt.Errorf("Command not supported: %v", c.args[0])
	}
	return nil
}

// generate walks the file looking for a handful of known tokens
func generate(c *command) error {
	//g := newGenerator()
	//g.Compare(l)
	// WalkFunc through and search each file. We do want to early exit on an unsupported mime
	// We'll use a codified list of mappings for now, PR to update.
	// use http.DetectContentType(data) and send off to a go routine when it's a known type
	// Return on channel to run l.command on for the given types, run commands in main thread until done
	// In the end, we'll have a workable layout
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
		l = &Layout{
			TaskList: []*Entries{},
		}
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
