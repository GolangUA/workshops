# Concurrent web crawler

Write a concurrent web crawler that concurrently executes http requests and is trying to find a word `concurrency` in the response data.

1. Create Worker function that perfroms requests and sends result into data channel.
2. Create Reader function that reads []bytes from the request channel and performs search in string of bytes.
3. Print which URL will contain the resulting data.
4. Your functions should have a cancelation-listening mechanism. If a result is found - stop all the goroutines and requests. (Use context and cancel func for that)


The output should be smth like:

```
starting sending request to https://google.com
starting sending request to https://twitter.com/
starting sending request to http://localhost:8000
starting sending request to https://itc.ua/
starting sending request to https://twitter.com/concurrencyinc
starting sending request to https://github.com/bradtraversy/go_restapi/blob/master/main.go
starting sending request to https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.
starting sending request to https://www.youtube.com/
starting sending request to https://postman-echo.com/get
Nothing found in https://github.com/bradtraversy/go_restapi/blob/master/main.go
'concurrency' string is found in https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.
Get "https://itc.ua/": context canceled
Get "https://twitter.com/": context canceled
Get "https://twitter.com/concurrencyinc": context canceled
Get "https://postman-echo.com/get": context canceled
exiting from searcher...
Get "https://www.google.com/": context canceled
Get "http://localhost:8000": context canceled

```