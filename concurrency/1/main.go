package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return tweets
		}
		// TODO: use channel here
		tweets = append(tweets, tweet)
	}
}

func consumer(tweets []*Tweet) {
	// TODO: use channel here
	for _, t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
			continue
		}

		fmt.Println(t.Username, "\tdoes not tweet about golang")
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Modification starts from here
	// Hint: this can be resolved via channels
	// Producer
	tweets := producer(stream)
	// Consumer
	consumer(tweets)

	fmt.Printf("Process took %s\n", time.Since(start))
}
