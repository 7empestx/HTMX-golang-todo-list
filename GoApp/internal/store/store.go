package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/7empestx/GoHTMXToDoList/internal/store/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
	q  *storedb.Queries
}

var dbInstance *Store

func InitDB(dataSourceName string) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	dbInstance = &Store{
		db: db,
		q:  storedb.New(db),
	}

	return nil
}

func GetTasks(ctx context.Context) ([]storedb.Task, error) {
	tasks, err := dbInstance.q.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var result []storedb.Task
	for _, task := range tasks {
		result = append(result, storedb.Task{
			ID:          task.ID,
			Description: task.Description,
			Completed:   task.Completed,
			Addedfrom:   task.Addedfrom,
		})
	}

	return result, nil
}

func AddTask(ctx context.Context, description string, addedFrom string) error {
	err := dbInstance.q.AddTask(ctx, storedb.AddTaskParams{
		Description: sql.NullString{String: description, Valid: true},
		Completed:   sql.NullBool{Bool: false, Valid: true},
		Addedfrom:   sql.NullString{String: addedFrom, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}

	return nil
}

func Checked(ctx context.Context, id int32) error {
	err := dbInstance.q.Checked(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check task: %w", err)
	}

	return nil
}

func DeleteTask(ctx context.Context, id int32) error {
	return dbInstance.q.DeleteTask(ctx, id)
}

func FilterCompletedTasks(ctx context.Context) ([]storedb.Task, error) {
	tasks, err := dbInstance.q.FilterCompletedTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to filter completed tasks: %w", err)
	}

	var result []storedb.Task
	for _, task := range tasks {
		result = append(result, storedb.Task{
			ID:          task.ID,
			Description: task.Description,
			Completed:   task.Completed,
			Addedfrom:   task.Addedfrom,
		})
	}

	return result, nil
}

func FilterIncompleteTasks(ctx context.Context) ([]storedb.Task, error) {
	tasks, err := dbInstance.q.FilterIncompleteTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to filter incomplete tasks: %w", err)
	}

	var result []storedb.Task
	for _, task := range tasks {
		result = append(result, storedb.Task{
			ID:          task.ID,
			Description: task.Description,
			Completed:   task.Completed,
			Addedfrom:   task.Addedfrom,
		})
	}

	return result, nil
}
