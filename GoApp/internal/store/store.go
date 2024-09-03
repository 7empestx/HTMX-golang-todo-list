package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/7empestx/GoHTMXToDoList/internal/store/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
	q  *storedb.Queries
}

var (
	dbInstance *Store
	once       sync.Once
	initErr    error
)

func InitDB(dataSourceName string) error {
	once.Do(func() {
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		dbInstance = &Store{
			db: db,
			q:  storedb.New(db),
		}
	})

	return initErr
}

func getStore() (*Store, error) {
	if dbInstance == nil {
		return nil, fmt.Errorf("database not initialized, call InitDB first")
	}
	return dbInstance, nil
}

func GetTasks(ctx context.Context) ([]storedb.Task, error) {
	store, err := getStore()
	if err != nil {
		return nil, err
	}

	tasks, err := store.q.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func AddTask(ctx context.Context, description string, addedFrom string) error {
	store, err := getStore()
	if err != nil {
		return err
	}

	err = store.q.AddTask(ctx, storedb.AddTaskParams{
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
	store, err := getStore()
	if err != nil {
		return err
	}

	err = store.q.Checked(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check task: %w", err)
	}
	return nil
}

func DeleteTask(ctx context.Context, id int32) error {
	store, err := getStore()
	if err != nil {
		return err
	}

	return store.q.DeleteTask(ctx, id)
}

func FilterCompletedTasks(ctx context.Context) ([]storedb.Task, error) {
	store, err := getStore()
	if err != nil {
		return nil, err
	}

	return store.q.FilterCompletedTasks(ctx)
}

func FilterIncompleteTasks(ctx context.Context) ([]storedb.Task, error) {
	store, err := getStore()
	if err != nil {
		return nil, err
	}

	return store.q.FilterIncompleteTasks(ctx)
}
