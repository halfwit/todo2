package main

import "fmt"

func task(c *command) error {
	// Open file
	// Parse in
	// Add/rm/etc
	// Write out
	switch c.args[1] {
	case "add":

	case "rm":
	case "do":
	case "undo":
	case "toggle":
	default:
		return fmt.Errorf("Command not supported: %v", c.args[1])
	}
	return nil
}
