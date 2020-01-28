package main

import "os"

type generator struct {
	existing *Layout
}

func (g *generator) dotTodoExists() bool {
	if _, err := os.Stat(".todo"); err != nil {
		return false
	}

	return true
}

func (g *generator) parseTodo() error {
	var err error
	
	g.existing, err = layoutFromTodoFile()
	if err != nil {
		return err
	}

	return nil
}

// WalkFunc through and search each file. We do want to early exit on an unsupported mime
// We'll use a codified list of mappings for now, PR to update.
// use http.DetectContentType(data) and send off to a go routine when it's a known type
// Return on channel to run l.command on for the given types, run commands in main thread until done
// In the end, we'll have a workable layout
// In our scrub function we need tags, so if there isn't one in a TODO entry, etc we will assume `[general]`
func (g *generator) toLayout() *Layout {
	// Use g.existing as a starting point and walk through our file path, add anything we can find.
	//go g.parseFile(filename)
	return nil
}
