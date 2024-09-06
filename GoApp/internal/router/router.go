package router

import (
  "log"
  "net/http"
  "github.com/7empestx/GoHTMXToDoList/internal/app"
  "github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()
  r.HandleFunc("/login", app.Login).Methods("GET", "POST")
  r.HandleFunc("/home", app.Home).Methods("GET")
	r.HandleFunc("/tasks", app.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", app.AddTask).Methods("POST")
	r.HandleFunc("/completed", app.FilterCompletedTasks).Methods("GET")
	r.HandleFunc("/incomplete", app.FilterIncompleteTasks).Methods("GET")
	r.HandleFunc("/checked", app.Checked).Methods("POST")
	r.HandleFunc("/delete/{id}", app.DeleteTask).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

  // Serve static files  
	log.Fatal(http.ListenAndServe(":8080", r))
  
}
