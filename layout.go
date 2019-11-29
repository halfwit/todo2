package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

/* This is what we want after parsing
type Layout struct {
	TaskList []*Tasks
}

type Tasks struct {
	Tag string
	Requires []string
	Tasks []*Task
}

type Task struct {
	Title string
	Entries []*Entry
}

type Entry struct {
	Desc string
	Done bool
}
*/

const (
	parent int = iota
	child
)

// Layout - Structure representing a given .todo file
type Layout struct {
	TaskList []*Entries
}

// Entries - set of todo items for a give task
type Entries struct {
	Title string
	List  []*Entry
}

// Entry - individual item of a set of tasks
type Entry struct {
	Desc string
	Done bool
}

func layoutFromWorkTree() (*Layout, error) {
	return nil, nil
}

func layoutFromTodoFile() (*Layout, error) {
	l := &Layout{
		TaskList: []*Entries{},
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

		tl := &Entries{
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

func parseEntries(sc *bufio.Scanner) []*Entry {
	sc.Scan()
	en := []*Entry{}
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
			en = append(en, &Entry{
				Desc: d[4:],
				Done: false,
			})
		case "[x]":
			en = append(en, &Entry{
				Desc: d[4:],
				Done: true,
			})
		default:
			return en
		}
		sc.Scan()
	}
}

func (l *Layout) destroy(title string) error {
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
	return errors.New("No such Entry")
}

func (l *Layout) create(title string) error {
	if len(title) < 1 {
		return errors.New("Unable to add nil Entry")
	}
	l.TaskList = append(l.TaskList, &Entries{
		Title: title,
		List:  []*Entry{},
	})
	return nil
}

func (l *Layout) taskExists(Title, item string) bool {
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

func (l *Layout) addTask(Title, item string) error {
	for _, e := range l.TaskList {
		if e.Title != Title {
			continue
		}
		e.List = append(e.List, &Entry{
			Desc: item,
			Done: false,
		})
		return nil
	}
	line := &Entry{
		Desc: item,
		Done: false,
	}
	l.TaskList = append(l.TaskList, &Entries{
		Title: Title,
		List:  []*Entry{line},
	})
	return nil
}

func (l *Layout) rmTask(Title, item string) error {
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
	return fmt.Errorf("No such task/Entry")
}

func (l *Layout) toggleTask(Title, item string) error {
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
	return fmt.Errorf("No such task/Entry")
}

func (l *Layout) addLink(n int, from, to string) {
	//for _, tasks := range l.TaskList {
	switch n {
	case parent:
		//if tasks.Tag == to {
		//	tasks.Requires = append(tasks.Requires, from)
		//}
	case child:
		//if tasks.Tag == from {
		//	tasks.Requires = append(tasks.Requires, to)
		//}
	}
	//}
}

func (l *Layout) rmLink(n int, from, to string) {
	//for _, tasks := range l.TaskList {
	switch n {
	case parent:
		//if tasks.Tag == to {
		//	for n, req := range tasks.Required {
		//		if req != from {
		//			continue
		//      }
		//		tasks.Required[n] = tasks.Required[len(tasks.Required)-1]
		//		tasks.Required = tasks.Required[:len(tasks.Required)-1]
		//	}
		//}
	case child:
		//if tasks.Tag == from {
		//	for n, req := range tasks.Required {
		//		if req != to {
		//			continue
		//      }
		//		tasks.Required[n] = tasks.Required[len(tasks.Required)-1]
		//		tasks.Required = tasks.Required[:len(tasks.Required)-1]
		//	}
		//}
	}
	//}
}
