package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)


const stringToSearch = "concurrency"

var sites = []string{
	"https://google.com",
	"https://itc.ua/",
	"https://twitter.com/concurrencyinc",
	"https://twitter.com/",
	"http://localhost:8000",
	"https://github.com/bradtraversy/go_restapi/blob/master/main.go",
	"https://www.youtube.com/",
	"https://postman-echo.com/get",
	"https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.",
}


type SiteData struct {
	data []byte
	uri  string
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	resultsCh := make(chan SiteData, len(sites))

	// TODO your code

	// give one second to validate if all other goroutines are closed
	time.Sleep(time.Second)
}

// TODO implement function that will execute request function, will validate the output and cancel all other requests when needed page is found
// and will listen to cancellation signal from context and will exit from the func when will receive it

// TODO implement function that will perfrom request using the example under

// TODO hint request function code:
/*
	Code to make request and read data

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

/*