package main

import (
	"errors"
	"fmt"
)

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
