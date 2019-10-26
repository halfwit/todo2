package main

import (
	"os"
	"text/template"
)

var tmpl = template.Must(template.New("new").Parse("{{.name}} test"))

type makefile struct {
	name string
}

func initFile(c *command) error {

	wr, err := os.Create(c.mkfile)
	if err != nil {
		return err
	}
	defer wr.Close()
	data := &makefile{"test"}
	return tmpl.Execute(wr, data)
}

func list(c *command) error {
	return nil
}
func listall(c *command) error {
	return nil
}
func dot(c *command) error {
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
