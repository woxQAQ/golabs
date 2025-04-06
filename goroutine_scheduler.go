package main

import (
	"fmt"
	"sync"
)

func GraphSchedule(m map[string][]string) {
	edges, verteies, vs := buildGraph(m)
	wg := &sync.WaitGroup{}
	wg.Add(len(vs))
	for k := range vs {
		go func() {
			processVertex(edges, verteies, k)
			wg.Done()
		}()
	}
	wg.Wait()
}

var processVertex = func(edges map[string][]chan struct{}, verteies map[string][]chan struct{}, vertex string) {
	for _, v := range edges[vertex] {
		<-v
	}
	fmt.Print(vertex)
	for _, v := range verteies[vertex] {
		v <- struct{}{}
	}
}

func buildGraph(m map[string][]string) (map[string][]chan struct{}, map[string][]chan struct{}, map[string]struct{}) {
	edegs := make(map[string][]chan struct{})
	verteies := make(map[string][]chan struct{})
	vs := make(map[string]struct{})
	for vertex, v := range m {
		vs[vertex] = struct{}{}
		for _, vv := range v {
			c := make(chan struct{})
			edegs[vertex] = append(edegs[vertex], c)
			verteies[vv] = append(verteies[vv], c)
			vs[vv] = struct{}{}
		}
	}
	return edegs, verteies, vs
}

// D: [B,C]
// B: [A]
// C: [A]
