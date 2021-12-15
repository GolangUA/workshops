# Dynamic worker pool

Write a dynamic worker pool that concurrently executes `tasks`. 

Worker pool should have `minimum` and `maximum` number of available workers 
When additional load comes - logic should increase number of available workers and when this load is decreased - remove them from pool.
Worker pool should also listen to cancellation signals. You should decide how the logic will understand if it still required to 
hold additional workers. Same goes to adding additional ones.

Hint - you can define time constraints. If there is no new load in 100 milliseconds - remove worker from pool.


Possible stuff to use: context, channels, atomics, selects, goroutines.


Regular worker pool simple example:
https://gobyexample.com/worker-pools