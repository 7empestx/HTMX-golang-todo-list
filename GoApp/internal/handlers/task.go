package handlers

import (
	"fmt"
	"net/http"
	"strconv"
  "context"

	"github.com/7empestx/GoHTMXToDoList/internal/models"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/7empestx/GoHTMXToDoList/internal/views"
	"github.com/gorilla/mux"
)

func renderComponent(w http.ResponseWriter, tasks []models.Task) {
  component := views.Tasks(tasks) 
  component.Render(context.Background(), w)
}

func FilterIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	tasks := store.FilterIncompleteTasks()
	renderComponent(w, tasks)
}

func FilterCompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks := store.FilterCompletedTasks()
	renderComponent(w, tasks)
}

func Checked(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the ID from the form data
	idStr := r.FormValue("taskID")

	if idStr == "" {
		fmt.Println("No taskID received")
		http.Error(w, "No taskID received", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error converting taskID to integer:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.Checked(id)
	GetTasks(w, r)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
  tasks := store.GetTasks()
  renderComponent(w, tasks)
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the description from the form data
	description := r.FormValue("description")

	var task models.Task
	task.Description = description

	store.AddTask(task.Description)
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
