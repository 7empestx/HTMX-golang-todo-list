package main

import (
	"log"
	"net/http"

	"github.com/7empestx/GoHTMXToDoList/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.AddTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
