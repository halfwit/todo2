package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"
)

const tmplt = `{{range .TaskList}}{{.Title}}
{{range .List}}{{if .Done}}[x] {{.Desc}}{{else}}[ ] {{.Desc}}{{end}}
{{end}}
{{end}}`

type layout struct {
	TaskList []*entries
}

// Entries - set of todo items for a give task
type entries struct {
	Title string
	List  []*entry
}

// Entry - individual item of a set of tasks
type entry struct {
	Desc string
	Done bool
}

func parse() (*layout, error) {
	l := &layout{
		TaskList: []*entries{},
	}
	fl, err := os.Open(".todo")
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(fl)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		tl := &entries{
			Title: sc.Text(),
			List:  parseEntries(sc),
		}
		l.TaskList = append(l.TaskList, tl)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return l, nil
}

func parseEntries(sc *bufio.Scanner) []*entry {
	sc.Scan()
	en := []*entry{}
	for {
		d := sc.Text()
		if len(d) < 4 {
			return en
		}
		if err := sc.Err(); err != nil {
			log.Fatal(err)
		}
		switch d[:3] {
		case "[ ]":
			en = append(en, &entry{
				Desc: d[4:],
				Done: false,
			})
		case "[x]":
			en = append(en, &entry{
				Desc: d[4:],
				Done: true,
			})
		default:
			return en
		}
		sc.Scan()
	}
}

func task(c *command) error {
	if len(c.args) < 2 {
		return errors.New("Too few arguments supplied")
	}
	l, err := parse()
	if err != nil && c.args[0] != "create" {
		return err
	}

	switch c.args[0] {
	case "create":
		l = &layout{
			TaskList: []*entries{},
		}
		err = l.create(c.args[1])
		if err != nil {
			return err
		}
		l.write()
		return nil
	case "destroy":
		err := l.destroy(c.args[1])
		if err != nil {
			return err
		}
		l.write()
		return nil
	case "add":
		if l.exists(c.args[1], c.args[2]) {
			return fmt.Errorf("Entry exists")
		}
		l.add(c.args[1], c.args[2])
		l.write()
		return nil
	case "rm":
		if !l.exists(c.args[1], c.args[2]) {
			return fmt.Errorf("No entry found")
		}
		defer l.write()
		return l.rm(c.args[1], c.args[2])
	case "toggle":
		if !l.exists(c.args[1], c.args[2]) {
			return fmt.Errorf("No such task/entry")
		}
		defer l.write()
		return l.toggle(c.args[1], c.args[2])
	default:
		return fmt.Errorf("Command not supported: %v", c.args[0])
	}
}

func (l *layout) destroy(title string) error {
	if _, err := os.Stat(".todo"); err != nil {
		return errors.New("Unable to locate .todo file")
	}
	for i, e := range l.TaskList {
		if e.Title != title {
			continue
		}
		if i < len(l.TaskList)-1 {
			copy(l.TaskList[i:], l.TaskList[i+1:])
		}
		l.TaskList = l.TaskList[:len(l.TaskList)-1]
		return nil
	}
	return errors.New("No such entry")
}

func (l *layout) create(title string) error {
	if len(title) < 1 {
		return errors.New("Unable to add nil entry")
	}
	l.TaskList = append(l.TaskList, &entries{
		Title: title,
		List:  []*entry{},
	})
	return nil
}

func (l *layout) write() {
	wr, err := os.Create(".todo")
	defer wr.Close()
	if err != nil {
		log.Fatal(err)
	}
	t := template.Must(template.New("tmplt").Parse(tmplt))
	err = t.Execute(wr, l)
	if err != nil {
		log.Fatal(err)
	}
}

func (l *layout) exists(Title, item string) bool {
	for _, e := range l.TaskList {
		if e.Title != Title {
			continue
		}
		for _, t := range e.List {
			if t.Desc == item {
				return true
			}
		}
	}
	return false
}

func (l *layout) add(Title, item string) error {
	for _, e := range l.TaskList {
		// Found existing task, append to List
		if e.Title != Title {
			continue
		}
		e.List = append(e.List, &entry{
			Desc: item,
			Done: false,
		})
		return nil
	}
	line := &entry{
		Desc: item,
		Done: false,
	}
	l.TaskList = append(l.TaskList, &entries{
		Title: Title,
		List:  []*entry{line},
	})
	return nil
}

func (l *layout) rm(Title, item string) error {
	for _, e := range l.TaskList {
		if e.Title != Title {
			continue
		}
		for i, t := range e.List {
			if t.Desc != item {
				continue
			}
			if i < len(e.List)-1 {
				copy(e.List[i:], e.List[i+1:])
			}
			e.List = e.List[:len(e.List)-1]
			return nil
		}
	}
	return fmt.Errorf("No such task/entry")
}

func (l *layout) toggle(Title, item string) error {
	for _, e := range l.TaskList {
		if e.Title != Title {
			continue
		}
		for _, t := range e.List {
			if t.Desc != item {
				continue
			}
			t.Done = !t.Done
			return nil
		}
	}
	return fmt.Errorf("No such task/entry")
}
