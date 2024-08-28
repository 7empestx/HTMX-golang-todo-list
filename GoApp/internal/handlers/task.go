package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/7empestx/GoHTMXToDoList/internal/store"
	"github.com/7empestx/GoHTMXToDoList/internal/views"
	"github.com/7empestx/GoHTMXToDoList/internal/login"
  "github.com/7empestx/GoHTMXToDoList/internal/store/sqlc"
	"github.com/gorilla/mux"
)

func Home(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Home handler called")
  component := login.Home()
  err := component.Render(r.Context(), w)
  if err != nil {
    fmt.Printf("Error rendering home component: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}

func Login(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Login handler called")
  component := login.Login()
  err := component.Render(r.Context(), w)
  if err != nil {
    fmt.Printf("Error rendering login component: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}

func renderComponent(w http.ResponseWriter, tasks []storedb.Task) {
	component := views.Tasks(tasks)
	component.Render(context.Background(), w)
}

func FilterIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := store.FilterIncompleteTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderComponent(w, tasks)
}

func FilterCompletedTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := store.FilterCompletedTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderComponent(w, tasks)
}

func Checked(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	ctx := context.Background()
	err = store.Checked(ctx, int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GetTasks(w, r)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := store.GetTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderComponent(w, tasks)
}

func logIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip := logIPAddress(r)
	description := r.FormValue("description")

	ctx := context.Background()
	err := store.AddTask(ctx, description, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GetTasks(w, r)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = store.DeleteTask(ctx, int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GetTasks(w, r)
}
