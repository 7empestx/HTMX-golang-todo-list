package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/7empestx/GoHTMXToDoList/internal/models"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tasks := store.GetTasks()
	var htmlResponse string
	for _, task := range tasks {
		htmlResponse += fmt.Sprintf("<li>%s</li>", task.Description)
	}
	fmt.Fprintf(w, htmlResponse)
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddTask")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the description from the form data
	description := r.FormValue("description")
	fmt.Println("Form Description:", description)

	var task models.Task
	task.Description = description

	fmt.Println("Decoded Task:", task)

	task = store.AddTask(task.Description)
	GetTasks(w, r) // Call GetTasks to render the updated list
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.DeleteTask(id)
	GetTasks(w, r) // Call GetTasks to render the updated list
}
