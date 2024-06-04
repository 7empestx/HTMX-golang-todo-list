package handlers

import (
	"fmt"
	"net/http"
	"strconv"
  "html/template"
  "log"

	"github.com/7empestx/GoHTMXToDoList/internal/models"
	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/gorilla/mux"
)

func Checked(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the ID from the form data
	idStr := r.FormValue("taskID")

  id, err := strconv.Atoi(idStr)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

	var task models.Task
	task.Completed = true

  store.Checked(id)
  GetTasks(w, r)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	const tmpl = `
    <li>
      {{if .Completed}}
        <input hx-post="/checked" hx-trigger="click" type="checkbox" id="taskID" name="taskID" value="{{ .ID }}" checked="true">
      {{else}}
        <input hx-post="/checked" hx-trigger="click" type="checkbox" id="taskID" name="taskID" value="{{ .ID }}">
      {{end}}
      <label for="taskID"> {{ .Description }}</label><br>
    </li>
  `

	w.Header().Set("Content-Type", "text/html")
	tasks := store.GetTasks()
  check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

  t, err := template.New("webpage").Parse(tmpl)
  check(err)

  for _, task := range tasks {
    if err := t.Execute(w, task); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
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
