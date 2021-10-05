Given is a producer-consumer scenario, where a producer reads in
tweets from a mockstream and a consumer is processing the
data. Your task is to change the code so that the producer as well
as the consumer can run concurrently both at the same time.

There are 2 files - mockstream.go and main.go \
You don't need to modify mockstream.go file. Your task is to modify code in main.go only.

# Producer-Consumer Scenario

The producer reads in tweets from a mockstream and a consumer is processing the data to find out whether someone has tweeted about golang or not. The task is to modify the code inside `main.go` so that producer and consumer can run concurrently to increase the throughput of this program.

## Expected results:
Before: 
```
davecheney      tweets about golang
beertocode      does not tweet about golang
ironzeb         tweets about golang
beertocode      tweets about golang
vampirewalk666  tweets about golang
Process took 3.580866005s
```

After:
```
davecheney      tweets about golang
beertocode      does not tweet about golang
ironzeb         tweets about golang
beertocode      tweets about golang
vampirewalk666  tweets about golang
Process took 1.977756255s
```