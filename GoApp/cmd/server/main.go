package main

import (
	"log"
	"net/http"

	"github.com/7empestx/GoHTMXToDoList/internal/handlers"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/gorilla/mux"
)

func main() {
	dataSourceName := "root:password@tcp(127.0.0.1:3306)/tododb"
	store.InitDB(dataSourceName)

	r := mux.NewRouter()
	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.AddTask).Methods("POST")
	r.HandleFunc("/completed", handlers.FilterCompletedTasks).Methods("GET")
	r.HandleFunc("/incomplete", handlers.FilterIncompleteTasks).Methods("GET")
	r.HandleFunc("/checked", handlers.Checked).Methods("POST")
	r.HandleFunc("/delete/{id}", handlers.DeleteTask).Methods("POST")

  // Serve static files  
	log.Fatal(http.ListenAndServe(":5000", r))
}
