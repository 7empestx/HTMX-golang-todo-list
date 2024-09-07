package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/7empestx/GoHTMXToDoList/internal/db/store/sqlc"
  "github.com/7empestx/GoHTMXToDoList/internal/db"
	_ "github.com/go-sql-driver/mysql"
)

func GetTasks(ctx context.Context) ([]storedb.Task, error) {
	store, err := db.GetStore()
	if err != nil {
		return nil, err
	}

	tasks, err := store.Q.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func AddTask(ctx context.Context, description string, addedFrom string) error {
	store, err := db.GetStore()
	if err != nil {
		return err
	}

	err = store.Q.AddTask(ctx, storedb.AddTaskParams{
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
	store, err := db.GetStore()
	if err != nil {
		return err
	}

	err = store.Q.Checked(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check task: %w", err)
	}
	return nil
}

func DeleteTask(ctx context.Context, id int32) error {
	store, err := db.GetStore()
	if err != nil {
		return err
	}

	return store.Q.DeleteTask(ctx, id)
}

func FilterCompletedTasks(ctx context.Context) ([]storedb.Task, error) {
	store, err := db.GetStore()
	if err != nil {
		return nil, err
	}

	return store.Q.FilterCompletedTasks(ctx)
}

func FilterIncompleteTasks(ctx context.Context) ([]storedb.Task, error) {
	store, err := db.GetStore()
	if err != nil {
		return nil, err
	}

	return store.Q.FilterIncompleteTasks(ctx)
}
