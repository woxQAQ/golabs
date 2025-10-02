package syncx

import "sync"

type TaskFunc func()

type Task struct {
	ID     string
	Func   TaskFunc
	Deps   []*Task
	mu     sync.Mutex
	done   bool
	doneCh chan struct{}
}

func NewTask(id string, fn TaskFunc) *Task {
	t := &Task{
		ID:     id,
		Func:   fn,
		doneCh: make(chan struct{}),
	}
	return t
}

func (t *Task) Done() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.done {
		t.done = true
		close(t.doneCh)
	}
}

func (t *Task) Wait() {
	<-t.doneCh
}
