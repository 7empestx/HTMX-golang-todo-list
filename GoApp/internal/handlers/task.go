package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/7empestx/GoHTMXToDoList/internal/models"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/gorilla/mux"
)

// Template for task list item
const taskTemplate = `
<li>
  <form hx-post="/checked" hx-trigger="click" hx-target="#task-list" hx-swap="innerHTML">
    {{if .Completed}}
      <input type="checkbox" id="taskID{{ .ID }}" name="taskID" value="{{ .ID }}" checked>
    {{else}}
      <input type="checkbox" id="taskID{{ .ID }}" name="taskID" value="{{ .ID }}">
    {{end}}
    <input type="hidden" name="taskID" value="{{ .ID }}">
    <label for="taskID{{ .ID }}"> {{ .Description }}</label><br>
  </form>
</li>
`

// Render tasks using the shared template
func renderTasks(w http.ResponseWriter, tasks []models.Task) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.New("webpage").Parse(taskTemplate)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		if err := t.Execute(w, task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func FilterIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	tasks := store.FilterIncompleteTasks()
	renderTasks(w, tasks)
}

func FilterCompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks := store.FilterCompletedTasks()
	renderTasks(w, tasks)
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
	renderTasks(w, tasks)
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
