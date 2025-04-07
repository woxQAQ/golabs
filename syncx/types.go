package syncx

import "sync"

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

func NewTask(Id string, Worker func()) Task {
	return Task{Id, Worker}
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
