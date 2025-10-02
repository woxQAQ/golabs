package syncx

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrCycleDetected = errors.New("cycle detected in dependency graph")
)

type Graph struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

func NewGraph() *Graph {
	return &Graph{
		tasks: make(map[string]*Task),
	}
}

func (g *Graph) AddTask(id string, fn TaskFunc) *Task {
	g.mu.Lock()
	defer g.mu.Unlock()

	task := NewTask(id, fn)
	g.tasks[id] = task
	return task
}

func (g *Graph) AddDependency(from, to string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	fromTask, exists := g.tasks[from]
	if !exists {
		return nil
	}

	toTask, exists := g.tasks[to]
	if !exists {
		return nil
	}

	if g.hasCycle(toTask, fromTask) {
		return ErrCycleDetected
	}

	toTask.Deps = append(toTask.Deps, fromTask)
	return nil
}

func (g *Graph) hasCycle(start, current *Task) bool {
	visited := make(map[*Task]bool)
	return g.dfsCycleCheck(start, current, visited)
}

func (g *Graph) dfsCycleCheck(start, current *Task, visited map[*Task]bool) bool {
	if current == start {
		return true
	}

	if visited[current] {
		return false
	}

	visited[current] = true

	for _, dep := range current.Deps {
		if g.dfsCycleCheck(start, dep, visited) {
			return true
		}
	}

	return false
}

func (g *Graph) GetTask(id string) *Task {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.tasks[id]
}

func (g *Graph) Execute() {
	g.ExecuteWithContext(context.Background())
}

func (g *Graph) ExecuteWithContext(ctx context.Context) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var wg sync.WaitGroup

	for _, task := range g.tasks {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			g.executeTask(ctx, t)
		}(task)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		return
	}
}

func (g *Graph) executeTask(ctx context.Context, task *Task) {
	select {
	case <-ctx.Done():
		task.Done()
		return
	default:
	}

	for _, dep := range task.Deps {
		select {
		case <-ctx.Done():
			task.Done()
			return
		default:
			dep.Wait()
		}
	}

	if task.Func != nil {
		task.Func()
	}

	task.Done()
}
