package store

import "github.com/7empestx/GoHTMXToDoList/internal/models"

var (
	tasks  = []models.Task{}
	nextID = 1
)

func FilterIncompleteTasks() []models.Task {
  var incompleteTasks []models.Task
  for _, task := range tasks {
    if !task.Completed {
      incompleteTasks = append(incompleteTasks, task)
    }
  }
  return incompleteTasks
}

func FilterCompletedTasks() []models.Task {
  var completedTasks []models.Task
  for _, task := range tasks {
    if task.Completed {
      completedTasks = append(completedTasks, task)
    }
  }
  return completedTasks
}

func Checked(id int) {
  for i, task := range tasks {
    if task.ID == id {
      tasks[i].Completed = !tasks[i].Completed
      break
    }
  }
}

func AddTask(description string) models.Task {
	task := models.Task{
		ID:          nextID,
		Description: description,
    Completed:   false,
	}
	tasks = append(tasks, task)
	nextID++
	return task
}

func DeleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
}

func GetTasks() []models.Task {
	return tasks
}

