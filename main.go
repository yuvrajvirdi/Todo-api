package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Todo struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "user:root@tcp(127.0.0.1:3306)/todos")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todos []Todo

	result, err := db.Query("SELECT id, task from todos")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var todo Todo
		err := result.Scan(&todo.ID, &todo.Task)
		if err != nil {
			panic(err.Error())
		}
		todos = append(todos, todo)
	}

	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {

	stmt, err := db.Prepare("INSERT INTO todos(task) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	task := keyVal["task"]

	_, err = stmt.Exec(task)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New Todo task added.")
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT id, task FROM todos WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var todo Todo

	for result.Next() {
		err := result.Scan(&todo.ID, &todo.Task)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE todos SET task = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTask := keyVal["task"]

	_, err = stmt.Exec(newTask, params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Changed Todo with ID %s", params["id"])
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Deleted Todo with ID %s", params["id"])
}
