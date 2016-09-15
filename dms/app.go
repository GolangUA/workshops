package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
}

type TaskList []Task

//Step3: Implement of interaction with database
type dbDriver interface {
	Create(t Task) error
	ReadById(id *int64) (TaskList, error)
	ReadByAlias(alias *string) (TaskList, error)
	Update(t Task) error
	Delete(t Task) error
}

//Step2: Create API to handles such type of calls or use exists routes
func handleProcess(w http.ResponseWriter, r *http.Request) {
	log.Printf("%d request %#v was processed\n", 0, nil)
}

//Step4: Implement CRUD handlers

//Step3: create connection with DB, docker-compose should be used for launch DB
func connect2Db() sql.DB {
}

func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/", handleProcess)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
