package store

import (
	"database/sql"
	"github.com/7empestx/GoHTMXToDoList/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTable()
	addColumnIfNotExists("tasks", "addedFrom", "VARCHAR(45)")
}

func createTable() {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        description TEXT,
        completed BOOLEAN
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func addColumnIfNotExists(tableName, columnName, columnType string) {
	var exists bool
	query := `
    SELECT COUNT(*) > 0
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = ?
    AND COLUMN_NAME = ?
    `
	err := db.QueryRow(query, tableName, columnName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		query = `
        ALTER TABLE ` + tableName + ` ADD ` + columnName + ` ` + columnType + `;
        `
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetTasks() []models.Task {
    rows, err := db.Query("SELECT id, description, completed FROM tasks")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        if err := rows.Scan(&task.ID, &task.Description, &task.Completed); err != nil {
            log.Fatal(err)
        }
        tasks = append(tasks, task)
    }

    return tasks
}

var (
	tasks  = []models.Task{}
	nextID = 1
)

func AddTask(description string, addedFrom string) {
	_, err := db.Exec("INSERT INTO tasks (description, completed, addedFrom) VALUES (?, ?, ?)", description, false, addedFrom)
	if err != nil {
		log.Fatal(err)
	}
}

func Checked(id int) {
	_, err := db.Exec("UPDATE tasks SET completed = NOT completed WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteTask(id int) {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

func FilterCompletedTasks() []models.Task {
	rows, err := db.Query("SELECT id, description, completed FROM tasks WHERE completed = TRUE")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Completed); err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func FilterIncompleteTasks() []models.Task {
	rows, err := db.Query("SELECT id, description, completed FROM tasks WHERE completed = FALSE")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Completed); err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}
