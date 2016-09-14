package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

//Step4: define the Message structure and add logic to the handler function
type Message struct {
}

//Step3: set global variable as message counter

//Step2: fill the API handler
func handleProcess(w http.ResponseWriter, r *http.Request) {
	log.Printf("%d Message %#v was processed\n", 0, nil)
}

//Step5: create connection with DB, docker-compose should be used for launch DB
func connect2Db() sql.DB {
}

func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/", handleProcess)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
