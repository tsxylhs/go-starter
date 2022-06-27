package task

import (
	"runtime"
	"sync"
	"time"
)

type Task struct {
	CountPool   int
	TimeOut     time.Time
	Results     chan interface{}
	Name        string
	Description []string
	Examples    []string
	Do          func() interface{}
}

// 1. 简单并发任务", "2.按时间来持续并发", "3.以 worker pool 方式 并发做事/发送请求", "4.等待异步任务执行结果"
func (task *Task) Task1() {
	var wg sync.WaitGroup
	for i := 0; i < task.CountPool; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			task.Do()
		}(i)
	}
	wg.Wait()
}

// 按照时间持续并发
func (task *Task) Task2() {

	n := runtime.NumCPU()
	waitForAll := make(chan struct{})
	done := make(chan struct{})
	concurrentCount := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		concurrentCount <- struct{}{}
	}
	go func() {
		for time.Now().Before(task.TimeOut) {
			<-done
			concurrentCount <- struct{}{}
		}
		waitForAll <- struct{}{}
	}()

	go func() {
		for {
			<-concurrentCount
			go func() {
				task.Do()
				done <- struct{}{}
			}()
		}
	}()
	<-waitForAll
}

// 以workerPool 方式并发做事
func (task *Task) Task3() {
	var wg sync.WaitGroup

	doFunc := func(result chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()
		res := task.Do()
		result <- res
	}
	for i := 0; i < task.CountPool; i++ {
		wg.Add(1)
		go doFunc(task.Results, &wg)
	}
	wg.Wait()
}
