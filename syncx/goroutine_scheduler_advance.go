package syncx

import (
	"sync"
)

func GraphScheduleV2(tasks []Task, deps map[string][]string) {
	graph := buildTask(tasks, deps)
	wg := &sync.WaitGroup{}
	wg.Add(len(graph))
	for _, v := range graph {
		v.run(wg)
	}
	wg.Wait()
}

func buildTask(tasks []Task, dependencies map[string][]string) map[string]*Vertex {
	res := make(map[string]*Vertex)
	for _, task := range tasks {
		curId := task.Id
		v := &Vertex{
			Task: task,
		}
		res[curId] = v
	}
	for to, deps := range dependencies {
		for _, from := range deps {
			e := edge{
				from:    res[to],
				to:      res[from],
				control: make(chan struct{}),
			}
			res[to].inEdges = append(res[to].inEdges, e)
			res[from].outEdges = append(res[from].outEdges, e)
		}
	}
	return res
}
