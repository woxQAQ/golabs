package syncx

import (
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	temp []string
	mu   sync.Mutex
)

func simpleWorker(Id string) func() {
	return func() {
		mu.Lock()
		defer mu.Unlock()
		temp = append(temp, Id)
	}
}

func TestScheduler(t *testing.T) {
	cases := []struct {
		name     string
		tasks    []Task
		deps     map[string][]string
		expected []string
	}{
		{
			"simple test",
			[]Task{
				NewTask("A", simpleWorker("A")),
				NewTask("B", simpleWorker("B")),
				NewTask("C", simpleWorker("C")),
				NewTask("D", simpleWorker("D")),
			},
			map[string][]string{
				"D": {"B", "C"},
				"B": {"A"},
				"C": {"A"},
			},
			[]string{
				"ABCD",
				"ACBD",
			},
		},
		{
			"complex test",
			[]Task{
				NewTask("A", simpleWorker("A")),
				NewTask("B", simpleWorker("B")),
				NewTask("C", simpleWorker("C")),
				NewTask("D", simpleWorker("D")),
				NewTask("E", simpleWorker("E")),
				NewTask("F", simpleWorker("F")),
			},
			map[string][]string{
				"D": {"B", "C"},
				"B": {"A"},
				"C": {"A"},
				"E": {"D"},
				"F": {"E"},
			},
			[]string{
				"ABCDEF",
				"ACBDEF",
				"ACDEBF",
				"ACDEBF",
			},
		},
	}
	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			GraphScheduleV2(cases[i].tasks, cases[i].deps)
			time.Sleep(1 * time.Second)
			if !slices.Contains(cases[i].expected, strings.Join(temp, "")) {
				t.Errorf("Expect %s, actual: %s", cases[i].expected, strings.Join(temp, ""))
			}
			temp = temp[:0]
		})
	}
}
