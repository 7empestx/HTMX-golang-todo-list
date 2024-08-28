package main

import (
  "fmt"
	"log"
	"net/http"
  "os"

	"github.com/7empestx/GoHTMXToDoList/internal/handlers"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/gorilla/mux"
)

func main() {

	dbHost := os.Getenv("RDS_HOSTNAME")
	dbName := os.Getenv("RDS_DB_NAME")
	dbUser := os.Getenv("RDS_USERNAME")
	dbPassword := os.Getenv("RDS_PASSWORD")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
  fmt.Println(dataSourceName)
	store.InitDB(dataSourceName)

	r := mux.NewRouter()
  r.HandleFunc("/login", handlers.Login).Methods("GET")
  r.HandleFunc("/home", handlers.Home).Methods("GET")
	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.AddTask).Methods("POST")
	r.HandleFunc("/completed", handlers.FilterCompletedTasks).Methods("GET")
	r.HandleFunc("/incomplete", handlers.FilterIncompleteTasks).Methods("GET")
	r.HandleFunc("/checked", handlers.Checked).Methods("POST")
	r.HandleFunc("/delete/{id}", handlers.DeleteTask).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

  // Serve static files  
	log.Fatal(http.ListenAndServe(":5000", r))
}
