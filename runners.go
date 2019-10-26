package main

import (
	"os"
	"text/template"
)

var tmpl = template.Must(template.New("new").Parse("# {{.name}} TODO"))

type makefile struct {
	name string
	//nodes []node
}

func initFile(c *command) error {

	wr, err := os.Create(c.mkfile)
	if err != nil {
		return err
	}
	defer wr.Close()
	//TODO(halfwit): Create a more useful type here
	data := &makefile{"test"}
	return tmpl.Execute(wr, data)
}

func list(c *command) error {
	// Parse makefile into DAG
	// Pretty print the leaves
	return nil
}
func listall(c *command) error {
	// Parse makefile into DAG
	// Pretty print everything
	return nil
}
func dot(c *command) error {
	// Parse makefile into DAG
	// Pretty print as dot format
	return nil
}
func rm(c *command) error {
	return nil
}
func add(c *command) error {
	return nil
}
func generate(c *command) error {
	return nil
}
