package main

import (
	"context"
	"sync/atomic"
)

type DynamicWP struct {
	// number of workers
	min, max           int
	currentWorkerCount *int32
}

// fill the correct arguments
func (w *DynamicWP) work(ctx context.Context, workerTasks chan func()) {
	atomic.AddInt32(w.currentWorkerCount, 1)

	// work should call the task function from task ch

	// don't forget to implement graceful shutdown:
	//case <-ctx.Done():
	//			atomic.AddInt32(w.currentWorkerCount, -1)
	//			return
	//}
}

// Start starts dynamic worker pull logic
func (w *DynamicWP) Start(ctx context.Context, tasksCh chan func()) {
	// worker load logic
}

func NewDynamicWorkerPool(min, max int) *DynamicWP {
	return &DynamicWP{
		min:                min,
		max:                max,
		currentWorkerCount: new(int32),
	}
}
