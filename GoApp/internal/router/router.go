package router

import (
  "log"
  "net/http"
  "github.com/7empestx/GoHTMXToDoList/internal/router/login"
  "github.com/7empestx/GoHTMXToDoList/internal/router/home"
  "github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()
  r.HandleFunc("/login", login.Login).Methods("GET", "POST")
  r.HandleFunc("/home", home.Home).Methods("GET")

	r.HandleFunc("/tasks", home.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", home.AddTask).Methods("POST")
	r.HandleFunc("/completed", home.FilterCompletedTasks).Methods("GET")
	r.HandleFunc("/incomplete", home.FilterIncompleteTasks).Methods("GET")
	r.HandleFunc("/checked", home.Checked).Methods("POST")
	r.HandleFunc("/delete/{id}", home.DeleteTask).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

  // Serve static files  
	log.Fatal(http.ListenAndServe(":8080", r))
  
}
