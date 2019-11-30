package main

import (
	"github.com/goombaio/dag"
)

func dagFromLayout(l *Layout) *dag.DAG {
	dg := dag.NewDAG()
	dm := make(map[string]*dag.Vertex)
	for _, job := range l.Jobs {
		dm[job.Key] = dag.NewVertex(job.Key, job)
		dg.AddVertex(dm[job.Key])
	}
	for _, job := range l.Jobs {
		for _, req := range job.Requires {
			for _, other := range l.Jobs {
				if tagsMatch(other.Tags, job.Tags) {
					continue
				}
				if !contains(other.Tags, req) {
					continue
				}
				dg.AddEdge(dm[job.Key], dm[other.Key])
			}
		}
	}
	return dg
}
