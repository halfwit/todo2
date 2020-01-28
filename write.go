package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/goombaio/dag"
)

// Templates are so fun.
const todoTmpl = `{{range .Jobs}}{{range $n, $t := .Tags}}{{if $n}} {{end}}[{{$t}}]{{end}}:{{ $length := len .Requires}}{{if gt $length 0}} requires{{range .Requires}} [{{.}}]{{end}}{{end}}
{{range .Tasks}}{{if .Title}}{{.Title}}
{{end}}{{range .Entries}}{{if .Done}}[x] {{.Desc}}
{{else}}[ ] {{.Desc}}
{{end}}{{end}}
{{end}}{{end}}`

// Writer that creates our normal file
func writeTodo(l *Layout) (string, error) {
	wr, err := os.Create(".todo")
	if err != nil {
		return "", err
	}

	defer wr.Close()

	t := template.Must(template.New("todoTmpl").Parse(todoTmpl))

	if e := t.Execute(wr, l); e != nil {
		return "", e
	}

	return "", nil
}

// Writer that outputs in dot format
func writeDot(l *Layout) (string, error) {
	var sb strings.Builder

	sb.WriteString("digraph depgraph {\n\trankdir=RL;\n")

	for _, job := range l.Jobs {
		sb.WriteString(fmt.Sprintf("%v", job.Key))
		sb.WriteString(` [label="`)

		switch len(job.Tasks) {
		case 1:
			for _, task := range job.Tasks {
				sb.WriteString(task.Title)

				for _, entry := range task.Entries {
					sb.WriteString("\n")
					sb.WriteString(entry.Desc)

					if entry.Done {
						sb.WriteString(" ✓")
					}
				}

				sb.WriteString("\n")
			}
		default:
			sb.WriteString(job.Tasks[0].Title)

			for _, entry := range job.Tasks[0].Entries {
				sb.WriteString("\n")
				sb.WriteString(entry.Desc)

				if entry.Done {
					sb.WriteString(" ✓")
				}
			}
		}

		sb.WriteString(`"];`)
		sb.WriteString("\r\n")
	}
	// Links
	for _, job := range l.Jobs {
		if len(job.Requires) == 0 {
			continue
		}

		for _, req := range job.Requires {
			sb.WriteString(fmt.Sprintf("%s -> %s;\n", req, job.Key))
		}
	}

	sb.WriteString("}\n")
	return sb.String(), nil
}

// Writer that outputs only leaves
func writeList(l *Layout) (string, error) {
	var sb strings.Builder

	l.removeCompleted()

	d, err := dagFromLayout(l)
	if err != nil {
		return "", err
	}

	leaves := d.SinkVertices()

	for _, leaf := range leaves {
		job := leaf.Value.(*Job)

		for _, t := range job.Tasks {
			fmt.Fprintf(&sb, "%v - %s\n", job.Tags, t.Title)

			for _, e := range t.Entries {
				var f rune

				switch e.Done {
				case true:
					f = '✓'
				case false:
					f = '✗'
				}

				fmt.Fprintf(&sb, " %c %s\n", f, e.Desc)
			}
		}
	}

	return sb.String(), nil
}

// Writer that outputs all nodes
func writeListAll(l *Layout) (string, error) {
	var walk func(v *dag.Vertex) error
	var sb strings.Builder

	d, err := dagFromLayout(l)
	if err != nil {
		return "", err
	}

	seen := map[*Job]bool{}

	walk = func(v *dag.Vertex) error {

		job := v.Value.(*Job)
		if seen[job] {
			return errors.New(errBadTodo)
		}

		seen[job] = true

		children, err := d.Successors(v)
		if err != nil {
			return err
		}

		for _, child := range children {
			walk(child)
		}

		for _, t := range job.Tasks {
			fmt.Fprintf(&sb, "[%v] - %s\n", job.Key, t.Title)

			for _, e := range t.Entries {
				var f rune

				switch e.Done {
				case true:
					f = '✓'
				case false:
					f = '✗'
				}

				fmt.Fprintf(&sb, " %c %s\n", f, e.Desc)
			}
		}

		return nil
	}

	for _, t := range d.SourceVertices() {
		if e := walk(t); e != nil {
			return "", e
		}
	}

	return sb.String(), nil
}
