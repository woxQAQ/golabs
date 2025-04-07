package syncx

import (
	"slices"
	"strings"
	"testing"
	"time"
)

var temp []string

func simpleWorker(Id string) func() {
	return func() {
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
			name: "",
			tasks: []Task{
				NewTask("A", simpleWorker("A")),
				NewTask("B", simpleWorker("B")),
				NewTask("C", simpleWorker("C")),
				NewTask("D", simpleWorker("D")),
			},
			deps: map[string][]string{
				"D": {"B", "C"},
				"B": {"A"},
				"C": {"A"},
			},
			expected: []string{
				"ABCD",
				"ACBD",
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
