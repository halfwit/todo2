package main

import (
	"fmt"
	"log"
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
func writeTodo(l *Layout) {
	wr, err := os.Create(".todo")
	defer wr.Close()
	if err != nil {
		log.Fatal(err)
	}
	t := template.Must(template.New("todoTmpl").Parse(todoTmpl))
	err = t.Execute(wr, l)
	if err != nil {
		log.Fatal(err)
	}
}

// Writer that outputs in dot format
func writeDot(l *Layout) {
	var sb strings.Builder
	sb.WriteString("digraph depgraph {\n\trankdir=RL;\n")
	// Labels
	for _, job := range l.Jobs {
		sb.WriteString(fmt.Sprintf("%v", job.Key))
		sb.WriteString(` [label="`)
		if len(job.Tasks) > 1 {
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
		} else {
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
			sb.WriteString(fmt.Sprintf("%s -> %s;\n", job.Key, req))
		}
	}
	sb.WriteString("}\n")
	fmt.Println(sb.String())
}

// Writer that outputs only leaves
func writeList(l *Layout) {
	l.removeCompleted()
	d := dagFromLayout(l)
	leaves := d.SinkVertices()
	for _, leaf := range leaves {
		job := leaf.Value.(*Job)
		for _, t := range job.Tasks {
			fmt.Printf("%v - %s\n", job.Tags, t.Title)
			for _, e := range t.Entries {
				var f rune
				switch e.Done {
				case true:
					f = '✓'
				case false:
					f = '✗'
				}
				fmt.Printf(" %c %s\n", f, e.Desc)
			}
		}
	}
}

// Writer that outputs all nodes
func writeListAll(l *Layout) {
	d := dagFromLayout(l)
	seen := map[*Job]bool{}
	var walk func(*dag.Vertex)
	walk = func(v *dag.Vertex) {
		job := v.Value.(*Job)
		if seen[job] {
			return
		}
		seen[job] = true
		children, err := d.Successors(v)
		if err != nil {
			log.Fatal(err)
		}
		for _, child := range children {
			walk(child)
		}
		for _, t := range job.Tasks {
			fmt.Printf("[%v] - %s\n", job.Key, t.Title)
			for _, e := range t.Entries {
				var f rune
				switch e.Done {
				case true:
					f = '✓'
				case false:
					f = '✗'
				}
				fmt.Printf(" %c %s\n", f, e.Desc)
			}
		}
	}
	for _, t := range d.SourceVertices() {
		walk(t)
	}
}
