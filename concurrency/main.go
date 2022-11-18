package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg            *sync.WaitGroup
	enlapsedTimes chan int64 = make(chan int64, 8)
)

func main() {
	wg = &sync.WaitGroup{}

	tasks := map[string]int64{
		"A": 123456789000,
		"B": 12345000,
		"C": 34547890000,
		"D": 100000000000,
		"E": 23234000,
		"F": 123000,
		"G": 897874000,
		"H": 9872343000,
	}

	fmt.Println("concurrency use case START at: " + time.Now().String())
	for task, count := range tasks {
		wg.Add(1)
		go func(task string, count int64, c chan int64) {
			c <- compute(task, int(count))
		}(task, count, enlapsedTimes)
	}

	wg.Wait()
	close(enlapsedTimes)
	fmt.Println("concurrency use case FINISH at: " + time.Now().String())

	totalTime := int64(0)
	for enlapsed := range enlapsedTimes {
		totalTime += enlapsed
	}
	fmt.Println("Time elapsed in milliseconds with concurrency: " + fmt.Sprint(totalTime))

	// single thread version
	fmt.Println("single thread use case START at: " + time.Now().String())

	totalTime = 0
	for task, count := range tasks {
		totalTime += compute2(task, int(count))
	}

	fmt.Println("single thread use case FINISH at: " + time.Now().String())
	fmt.Println("Time elapsed in milliseconds single thread: " + fmt.Sprint(totalTime))
}

func compute(taskName string, count int) int64 {
	defer wg.Done()

	startAt := time.Now()
	for x := 1; x <= count; x++ {
		_ = x * x
	}

	fmt.Println("finish task: " + taskName)
	return time.Since(startAt).Milliseconds()
}

func compute2(taskName string, count int) int64 {
	startAt := time.Now()
	for x := 1; x <= count; x++ {
		_ = x * x
	}

	fmt.Println("finish task: " + taskName)
	return time.Since(startAt).Milliseconds()
}
