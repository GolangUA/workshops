package main

import (
	"context"
	"sync/atomic"
	"time"
)

type DynamicWP struct {
	// number of workers
	min, max int
}

// fill the correct arguments
func (w *DynamicWP) work(ctx context.Context, tasksCh chan func()) {
	// work should call the task function from task ch
}

func (w *DynamicWP) Start(ctx context.Context, tasksCh chan func()) {
	// worker load logic
}

func NewDynamicWorkerPool(min, max int) *DynamicWP {
	return &DynamicWP{
		min: min,
		max: max,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	wp := NewDynamicWorkerPool(3, 20)

	tasksCh := make(chan func())
	go wp.Start(ctx, tasksCh)

	for i := 0; i < 100; i++ {
		tasksCh <- TestLoad
	}

	cancel() // worker pool should be gracefully closed
	time.Sleep(time.Second * 1)
}

var loadCount int32

func TestLoad() {
	atomic.AddInt32(&loadCount, 1)
	time.Sleep(time.Second * 1)
}
