# Dynamic worker pool

Write a dynamic worker pool that concurrently executes `tasks`. 

###Requirements:
1. Worker pool should have `minimum` and `maximum` number of available workers 
2. When additional load comes - logic should increase number of available workers and when this load is decreased - remove them from 'pool' by exiting from the goroutine.
3. Worker exit criteria - no load within 2 seconds.
4. Worker pool should also listen to cancellation signals and handle graceful shutdown. 
5. You should decide how the logic will understand if it still required to 
hold additional workers. 
6. New workers should be added if taskSame goes to adding additional ones.
7. New worker is added if task is not picked up within 50 milliseconds.


Possible stuff to use: context, channels, atomics, selects, goroutines.


Regular worker pool simple example:
https://gobyexample.com/worker-pools


How to test:
run 

    go mod download
    go test ./...