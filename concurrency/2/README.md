# Limit Service Time for Free-tier Users

Your video processing service has a freemium model. Every user has 10
sec of free processing time on your service. After that, the
service will kill your process, unless you are a paid premium user.
By `Kill` i mean that `HandleReuqest` function will return `false` value and that's all. 

You need to modify main.go file only.

To verify - just run code from main package and you will see the output.

Beginner Level Task: 10s max per request.
Advanced Level Task: 10s max per user (accumulated). So a user can have only 10 seconds for a request at all. For example
if 1st request was performed in 4 seconds, same user will have 6 seconds left to perform next request.

