package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Layout represents a fully parsed .todo file
type Layout struct {
	Jobs []*Job
}

// Job - A tagged collection of tasks
type Job struct {
	Key      string
	Tags     []string
	Requires []string
	Tasks    []*Task
}

// Task - A collection of Entries (checkboxes)
type Task struct {
	Title   string
	Entries []*Entry
}

// Entry -  single checkbox todo item
type Entry struct {
	Desc string
	Done bool
}

func layoutFromTodoFile() (*Layout, error) {

	l := &Layout{
		Jobs: []*Job{},
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
		// Build out entries in a stepwise-fashion. A proper parser grammer would be better here
		// but this is quick and assumes machine-written input
		j := &Job{}
		j.Tags = parseTags(sc)
		j.Key = strings.Join(j.Tags, ", ")
		j.Requires = parseRequires(sc)
		j.Tasks = parseTasks(sc)
		l.Jobs = append(l.Jobs, j)

	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return l, nil
}

func parseTags(sc *bufio.Scanner) []string {
	txt := sc.Text()
	n := strings.Index(txt, ":")
	return stringToTags(txt[:n])
}

func parseRequires(sc *bufio.Scanner) []string {
	txt := sc.Text()
	defer sc.Scan()
	n := strings.Index(txt, ": requires ")
	if len(txt) > n+len(": requires [ ]") {
		return stringToTags(txt[n+len(": requires "):])
	}
	return nil
}

func parseTasks(sc *bufio.Scanner) []*Task {
	var tasks []*Task
	for {
		title := sc.Text()
		if len(sc.Text()) == 0 {
			return tasks
		}
		if strings.HasPrefix(title, "[ ]") || strings.HasPrefix(title, "[x]") {
			title = ""
		} else {
			sc.Scan()
		}
		t := &Task{
			Title:   title,
			Entries: findEntries(sc),
		}
		tasks = append(tasks, t)
	}
}

func findEntries(sc *bufio.Scanner) []*Entry {
	en := []*Entry{}
	for {
		d := sc.Text()
		if len(d) < 4 {
			return en
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

func contains(first []string, last string) bool {
	for _, f := range first {
		if f == last {
			return true
		}
	}
	return false
}

func tagsMatch(first, last []string) bool {
	if len(first) != len(last) {
		return false
	}

	for n, f := range first {
		if f != last[n] {
			return false
		}
	}
	return true
}

func isCompleted(job *Job) bool {
	for _, task := range job.Tasks {
		for _, entry := range task.Entries {
			if !entry.Done {
				return false
			}
		}
	}
	return true
}

// This is adequate for our needs.
func stringToJobs(incoming string) (*Job, error) {
	r := regexp.MustCompile(`([^\[\]]*)\s?(\[.*\])+\s?(.*)`)
	matches := r.FindAllStringSubmatch(incoming, -1)
	if matches != nil && len(matches[0]) < 3 {
		return nil, errors.New("Could not parse string")
	}

	title := matches[0][1]
	if len(matches[0]) == 4 {
		title = fmt.Sprintf("%s%s", matches[0][1], matches[0][3])
	}
	t := &Task{
		Title: title,
	}
	return &Job{
		Tags:  stringToTags(matches[0][2]),
		Tasks: []*Task{t},
	}, nil
}

func trimBrackets(incoming string) string {
	incoming = strings.TrimPrefix(incoming, "[")
	return strings.TrimSuffix(incoming, "]")
}

func stringToTags(incoming string) []string {
	var tags []string
	var tmp string
	var startTag bool
	n := 0
	for _, b := range []byte(incoming) {
		if b == '[' {
			startTag = true
			tmp = ""
			continue
		}
		if b == ']' {
			n++
			tags = append(tags, tmp)
			startTag = false
		}
		if startTag {
			tmp += string(b)
		}
	}
	if tags == nil || tags[0] == "" {
		return nil
	}
	return tags
}

func (l *Layout) removeCompleted() {
	for i, job := range l.Jobs {
		if !isCompleted(job) {
			continue
		}
		if i < len(l.Jobs)-1 {
			copy(l.Jobs[i:], l.Jobs[i+1:])
		}
		l.Jobs = l.Jobs[:len(l.Jobs)-1]
	}
}

func (l *Layout) destroy(title string) error {
	if _, err := os.Stat(".todo"); err != nil {
		return errors.New("Unable to locate .todo file")
	}
	tasks, err := stringToJobs(title)
	if err != nil {
		return err
	}
	for i, e := range l.Jobs {
		if tagsMatch(e.Tags, tasks.Tags) {
			continue
		}
		if i < len(l.Jobs)-1 {
			copy(l.Jobs[i:], l.Jobs[i+1:])
		}
		l.Jobs = l.Jobs[:len(l.Jobs)-1]
		return nil

	}
	return errors.New("No such Entry")
}

func (l *Layout) create(incoming string) error {
	tasks, err := stringToJobs(incoming)
	if err != nil {
		return err
	}
	l.Jobs = append(l.Jobs, tasks)
	return nil
}

func (l *Layout) taskExists(title, item string) bool {
	t, err := stringToJobs(title)
	if err != nil {
		return false
	}
	for _, e := range l.Jobs {
		if !tagsMatch(e.Tags, t.Tags) {
			continue
		}
		for _, t := range e.Tasks {
			for _, d := range t.Entries {
				if d.Desc == item {
					return true
				}
			}
		}
	}
	return false
}

// TODO(halfwit) Verify that task title doesn't already exist!
func (l *Layout) addTask(title, item string) error {
	t, err := stringToJobs(title)
	if err != nil {
		return err
	}
	entry := &Entry{
		Desc: item,
		Done: false,
	}
	for _, job := range l.Jobs {
		if !tagsMatch(job.Tags, t.Tags) {
			continue
		}
		for _, task := range job.Tasks {
			if task.Title == t.Tasks[0].Title {
				task.Entries = append(task.Entries, entry)
				return nil
			}
		}
	}
	t.Tasks[0].Entries = append(t.Tasks[0].Entries, entry)
	for _, job := range l.Jobs {
		if !tagsMatch(job.Tags, t.Tags) {
			continue
		}
		job.Tasks = append(job.Tasks, t.Tasks[0])
		return nil
	}
	l.Jobs = append(l.Jobs, t)
	return nil
}

func (l *Layout) rmTask(title, item string) error {
	t, err := stringToJobs(title)
	if err != nil {
		return err
	}
	for _, e := range l.Jobs {
		if !tagsMatch(e.Tags, t.Tags) {
			continue
		}
		for _, t := range e.Tasks {
			for i, j := range t.Entries {
				if j.Desc != item {
					continue
				}
				if i < len(t.Entries)-1 {
					copy(t.Entries[i:], t.Entries[i+1:])
				}
				t.Entries = t.Entries[:len(t.Entries)-1]
				return nil
			}
		}
	}
	return fmt.Errorf("No such task/Entry")
}

func (l *Layout) toggleTask(title, item string) error {
	t, err := stringToJobs(title)
	if err != nil {
		return err
	}
	for _, e := range l.Jobs {
		if !tagsMatch(e.Tags, t.Tags) {
			continue
		}
		for _, t := range e.Tasks {
			for _, j := range t.Entries {
				if j.Desc != item {
					continue
				}
				j.Done = !j.Done
				return nil
			}
		}
	}
	return fmt.Errorf("No such task/Entry")
}

func (l *Layout) addLink(to, from string) {
	to = trimBrackets(to)
	from = trimBrackets(from)
	for _, tasks := range l.Jobs {
		for _, tag := range tasks.Tags {
			if tag != to {
				continue
			}
			if contains(tasks.Requires, from) {
				continue
			}
			tasks.Requires = append(tasks.Requires, from)
		}
	}
}

// It's ugly, but it gets the thing done
func (l *Layout) rmLink(to, from string) {
	for _, tasks := range l.Jobs {
		for _, tag := range tasks.Tags {
			if tag != to {
				continue
			}
			for n, req := range tasks.Requires {
				if req != from {
					continue
				}
				tasks.Requires[n] = tasks.Requires[len(tasks.Requires)-1]
				tasks.Requires = tasks.Requires[:len(tasks.Requires)-1]
			}

		}
	}
}
