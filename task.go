package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

const tmplt = `{{range .tasklist .}}
{{if and .list .title}}
{{.title}}
{{range .list .}}{{if .done}}[*]{{else}}[ ]{{end}} {{.desc}}{{end}}

{{end}}
{{end}}
`

/* Example:
My entry title
[ ] an entry
[*] a done entry

My other entry
[ ] an entry
[*] a done entry
*/

type layout struct {
	tasklist []entries
}

type entries struct {
	title string
	list  []entry
}

type entry struct {
	desc string
	done bool
}

func task(c *command) error {
	l, err := parse("test")
	if err != nil {
		return err
	}
	switch c.args[0] {
	case "add":
		if len(c.args) < 3 || l.exists(c.args[1], c.args[2]) {
			return fmt.Errorf("Entry exists")
		}
		defer l.write()
		return l.add(c.args[1], c.args[2])
	case "rm":
		if len(c.args) < 3 || !l.exists(c.args[1], c.args[2]) {
			return fmt.Errorf("No entry found")
		}
		defer l.write()
		return l.rm(c.args[1], c.args[2])
	case "toggle":
		if len(c.args) < 3 || !l.exists(c.args[1], c.args[2]) {
			return fmt.Errorf("No such task/entry")
		}
		defer l.write()
		return l.toggle(c.args[1], c.args[2])
	default:
		return fmt.Errorf("Command not supported: %v", c.args[0])
	}
}

func (l *layout) write() {
	wr, err := os.Create("tmp")
	defer wr.Close()
	if err != nil {
		log.Fatal(err)
	}
	t := template.Must(template.New(".todo").Parse(tmplt))
	t.Execute(wr, l)
}

func (l *layout) exists(title, item string) bool {
	for _, e := range l.tasklist {
		if e.title != title {
			continue
		}
		for _, t := range e.list {
			if t.desc == item {
				return true
			}
		}
	}
	return false
}

func (l *layout) add(title, item string) error {
	for _, e := range l.tasklist {
		if e.title != title {
			continue
		}
		e.list = append(e.list, entry{
			desc: item,
			done: false,
		})
		return nil
	}
	line := entry{
		desc: item,
		done: false,
	}
	l.tasklist = append(l.tasklist, entries{
		title: title,
		list:  []entry{line},
	})
	return nil
}

func (l *layout) rm(title, item string) error {
	for _, e := range l.tasklist {
		if e.title != title {
			continue
		}
		for i, t := range e.list {
			if t.desc != item {
				continue
			}
			if i < len(e.list)-1 {
				copy(e.list[i:], e.list[i+1:])
			}
			e.list = e.list[:len(e.list)-1]
			return nil
		}
	}
	return fmt.Errorf("No such task/entry")
}
func (l *layout) toggle(title, item string) error {
	for _, e := range l.tasklist {
		if e.title != title {
			continue
		}
		for _, t := range e.list {
			if t.desc != item {
				continue
			}
			t.done = !t.done
			return nil
		}
	}
	return fmt.Errorf("No such task/entry")
}
