package main

// WalkFunc through and search each file. We do want to early exit on an unsupported mime
// We'll use a codified list of mappings for now, PR to update.
// use http.DetectContentType(data) and send off to a go routine when it's a known type
// Return on channel to run l.command on for the given types, run commands in main thread until done
// In the end, we'll have a workable layout

type generator struct {
}

func newGenerator() *generator {
	return nil
}

func (g *generator) dotTodoExists() bool {
	return false
}

func (g *generator) parseTodo() {

}

// In our scrub function we need tags, so if there isn't one in a TODO entry, etc we will assume `[general]`
func (g *generator) toLayout() *Layout {
	return nil
}
