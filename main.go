package main

import (
	"fmt"
	"golabs/syncx"
)

func main() {
	input := map[string][]string{
		"D": {"B", "C"},
		"B": {"A"},
		"C": {"A"},
	}
	t := []syncx.Task{
		{
			Id:     "A",
			Worker: func() { fmt.Print("A") },
		},
		{
			Id:     "B",
			Worker: func() { fmt.Print("B") },
		},
		{
			Id:     "C",
			Worker: func() { fmt.Print("C") },
		},
		{
			Id:     "D",
			Worker: func() { fmt.Print("D") },
		},
	}
	syncx.GraphScheduleV2(t, input)
}
