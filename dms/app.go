package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type Message struct {
	EventType string
	Ts        int64
	Params    map[string]interface{}
}

var messCnt int64

func handleProcess(w http.ResponseWriter, r *http.Request) {
	messCnt := atomic.AddInt64(&messCnt, 1)
	log.Printf("Message %d was processed\n", messCnt)
}

func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/", handleProcess)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
