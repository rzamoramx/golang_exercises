package main

import (
	"fmt"
	"time"
)

var startAt time.Time

func main() {
	tasks := map[string]int64{
		"A": 123456789,
		"B": 12345,
		"C": 34547890,
		"D": 100000000,
		"E": 23234,
		"F": 123,
		"G": 897874,
		"H": 9872343,
	}

	startAt = time.Now()

	for task, count := range tasks {
		go compute(task, int(count))
	}

	enlapsed := time.Since(startAt)

	fmt.Println("Time: " + enlapsed.String())
}

func compute(taskName string, count int) {
	fmt.Println("entering task: " + taskName)

	for x := 1; x <= count; x++ {
		_ = x * x
	}

	fmt.Println("finish task")
}
