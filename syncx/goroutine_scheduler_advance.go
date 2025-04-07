package syncx

import (
	"sync"
)

type Task struct {
	Id     string
	Worker func()
}

type edge struct {
	from, to *Vertex
	control  chan struct{}
}

type Vertex struct {
	inEdges  []edge
	outEdges []edge
	Task
}

func (t *Vertex) run(ctrl *sync.WaitGroup) {
	go func() {
		for _, v := range t.inEdges {
			<-v.control
		}
		t.Worker()
		for _, v := range t.outEdges {
			v.control <- struct{}{}
		}
		ctrl.Done()
	}()
}

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
