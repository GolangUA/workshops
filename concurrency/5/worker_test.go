package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewDynamicWorkerPool(t *testing.T) {
	var (
		ctx, cancel     = context.WithCancel(context.TODO())
		maxWorkersCount = 20
		minWorkersCount = 3
		tasksCh         = make(chan func(), 100)
		taskFunc        = func() { time.Sleep(time.Second * 1) }
		wp              = NewDynamicWorkerPool(minWorkersCount, maxWorkersCount)
	)

	go wp.Start(ctx, tasksCh)
	for i := 0; i < 100; i++ {
		select {
		case tasksCh <- taskFunc:
		case <-time.After(time.Second * 4):
			t.Fatalf("task is not executed after 4 seconds")
		}
	}

	//---------------
	fmt.Println("Testing 'max load test case'")
	time.Sleep(time.Second * 2)
	currentWorkers := atomic.LoadInt32(wp.currentWorkerCount)
	require.EqualValues(t, maxWorkersCount, currentWorkers,
		"there should be max number of workers %v, but got %v", maxWorkersCount, currentWorkers)
	fmt.Println("Successfully passed")
	//---------------
	fmt.Println("Testing 'load has gone' test case")
	time.Sleep(time.Second * 10)
	currentWorkers = atomic.LoadInt32(wp.currentWorkerCount)
	require.EqualValues(t, minWorkersCount, currentWorkers,
		"workers should have been decreased when load is over to min number '%v', but got '%v'", minWorkersCount, currentWorkers)
	fmt.Println("Successfully passed")
	//---------------
	fmt.Println("Testing 'max load is back again' test case")
	for i := 0; i < 100; i++ {
		tasksCh <- taskFunc
	}
	time.Sleep(time.Second * 2)
	currentWorkers = atomic.LoadInt32(wp.currentWorkerCount)
	require.EqualValues(t, maxWorkersCount, currentWorkers,
		"workers count should have been maxed. expected '%v' workers, but got '%v'", maxWorkersCount, currentWorkers)
	fmt.Println("Successfully passed")
	//---------------
	fmt.Println("Testing 'decreased and canceled' test case")
	cancel() // worker pool should be gracefully closed
	time.Sleep(time.Second * 10)

	currentWorkers = atomic.LoadInt32(wp.currentWorkerCount)
	require.EqualValues(t, 0, currentWorkers,
		"workers should have been all closed when context is canceled, but got '%v' number of workers alive", currentWorkers)
	fmt.Println("Successfully passed")
}
