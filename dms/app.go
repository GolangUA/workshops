package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

type Message struct {
	EventType string                 `json:"eventType"`
	Ts        int64                  `json:"ts"`
	Params    map[string]interface{} `json:"params"`
}

var messCnt int64

func handleProcess(w http.ResponseWriter, r *http.Request) {
	var ms Message
	err := json.NewDecoder(r.Body).Decode(&ms)
	if err != nil {
		log.Printf("Can't decode JSON: %v", err)
		return
	}
	messCnt := atomic.AddInt64(&messCnt, 1)
	log.Printf("%d Message %#v was processed\n", messCnt, ms)
}

func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/", handleProcess)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
