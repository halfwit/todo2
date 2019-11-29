package main

import (
	"log"
	"os"
	"text/template"
)

// Templates are so fun.
const todoFile = `{{range .Jobs}}{{range $n, $t := .Tags}}{{if $n}} {{end}}[{{$t}}]{{end}}:{{ $length := len .Requires}}{{if gt $length 0}} requires{{range .Requires}} [{{.}}]{{end}}{{end}}
{{range .Tasks}}{{if .Title}}{{.Title}}
{{end}}{{range .Entries}}{{if .Done}}[x] {{.Desc}}{{else}}[ ] {{.Desc}}{{end}}
{{end}}
{{end}}{{end}}`

// Writer that creates our normal file
func writeTodo(l *Layout) {
	wr, err := os.Create(".todo")
	defer wr.Close()
	if err != nil {
		log.Fatal(err)
	}
	t := template.Must(template.New("todoFile").Parse(todoFile))
	err = t.Execute(wr, l)
	if err != nil {
		log.Fatal(err)
	}
}

// Writer that outputs in dot format
func writeDot(l *Layout) {

}

// Writer that outputs only leaves
func writeList(l *Layout) {

}

// Writer that outputs all nodes
func writeListAll(l *Layout) {

}
