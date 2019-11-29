package main

import (
	"log"
	"os"
	"text/template"
)

const tmplt = `{{range .TaskList}}{{.Title}}
{{range .List}}{{if .Done}}[x] {{.Desc}}{{else}}[ ] {{.Desc}}{{end}}
{{end}}
{{end}}`

// Writer that creates our normal file
func writeTodo(l *Layout) {
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

// Writer that outputs in dot format

// Writer that outputs only leaves

// Writer that outputs all nodes
