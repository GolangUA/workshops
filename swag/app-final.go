//go:generate swagger generate spec
// Package classification Petstore API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v2
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: John Doe<john.doe@example.com> http://john.doe.com
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
//
// swagger:meta
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID        int64    `json:"id,omitempty"`
	Alias     string   `json:"alias"`
	Desc      string   `json:"desc"`
	Category  []string `json:"cat,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Ts        int64    `json:"ts"`
	EstTime   string   `json:"est_time"`
	RealTime  string   `json:"real_time"`
	Reminders []string `json:"reminders,omitempty"`
}

type TaskList []Task

//Step3: Implement of interaction with database
type dbDriver interface {
	init() error
	Create(t Task) error
	read(v interface{}) (TaskList, error)
	ReadById(id *int64) (TaskList, error)
	ReadByAlias(alias *string) (TaskList, error)
	Update(t Task) error
	Delete(t Task) error
}

type App struct {
	st dbDriver
}

//Step2: Create API to handles such type of calls or use exists routes
func (a *App) handleProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		a.Read(w, r)
	case http.MethodPost:
		a.Create(w, r)
	case http.MethodPut:
		a.Update(w, r)
	case http.MethodDelete:
		a.Delete(w, r)
	}
	log.Printf("%d request %#v was processed\n", 0, nil)
}

//Step4: Implement CRUD handlers
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	var t Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Printf("Can't decode JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.st.Create(t)
	if err != nil {
		log.Printf("Can't create new task: %v (%#v)\n", err, t)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) Read(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[1:]
	var tl TaskList
	var err error
	var id int64
	if len(param) == 0 {
		tl, err = a.st.read(nil)
	} else if strings.IndexRune(param, '/') > -1 {
		log.Printf("URL contains more than one parameter: %s\n", param)
		http.Error(w, "URL contains more than one parameter", http.StatusBadRequest)
		return
	} else if id, err = strconv.ParseInt(param, 10, 64); err == nil {
		tl, err = a.st.ReadById(&id)
	} else {
		tl, err = a.st.ReadByAlias(&param)
	}
	if err != nil {
		log.Printf("Some error in select: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var js []byte
	js, err = json.Marshal(tl)
	if err != nil {
		log.Printf("Couldn't marshal list of tasks: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)

}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int64
	if id, err = strconv.ParseInt(r.URL.Path[1:], 10, 64); err != nil {
		log.Printf("URL didn't contains ID as parameter: %s\n", r.URL.Path)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var t Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Printf("Can't decode JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id != t.ID {
		log.Printf("ID from URL and JSON are different: %d <-> %d\n", id, t.ID)
		http.Error(w, "ID not match", http.StatusBadRequest)
		return
	}
	err = a.st.Update(t)
	if err != nil {
		log.Printf("Error while update of task: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	var id int64
	var err error
	if id, err = strconv.ParseInt(r.URL.Path[1:], 10, 64); err != nil {
		log.Printf("URL didn't contains ID as parameter: %s\n", r.URL.Path)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t := Task{ID: id}
	err = a.st.Delete(t)
	if err != nil {
		log.Printf("Can't delete the Task(%d): %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	log.Printf("Starting on port 8080")
	stDr := &sqliteDr{}
	err := stDr.init()
	if err != nil {
		log.Fatalf("can not connect to DB: %v", err)
	}
	a := &App{st: stDr}
	http.HandleFunc("/", a.handleProcess)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
