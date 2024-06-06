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
  r.HandleFunc("/checked", handlers.Checked).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	log.Fatal(http.ListenAndServe(":5000", r))
}
