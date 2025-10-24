package db

import (
  "context"
  "fmt"
  "os"
  "database/sql"
  "github.com/7empestx/GoHTMXToDoList/internal/db/store/sqlc"
  "github.com/7empestx/GoHTMXToDoList/internal/db/store/memory"
	_ "github.com/go-sql-driver/mysql"
)

// Querier defines the interface for database operations
type Querier interface {
	GetTasks(ctx context.Context) ([]storedb.Task, error)
	AddTask(ctx context.Context, arg storedb.AddTaskParams) error
	Checked(ctx context.Context, id int32) error
	DeleteTask(ctx context.Context, id int32) error
	FilterCompletedTasks(ctx context.Context) ([]storedb.Task, error)
	FilterIncompleteTasks(ctx context.Context) ([]storedb.Task, error)
}

type Store struct {
	db *sql.DB
	Q  Querier
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
  fmt.Println("Connected to database")

	dbInstance = &Store{
		db: db,
		Q:  storedb.New(db),
	}
  fmt.Println("Store initialized")

	return nil
}

func Init() error {
  dbHost := os.Getenv("RDS_HOSTNAME")
  dbName := os.Getenv("RDS_DB_NAME")
  dbUser := os.Getenv("RDS_USERNAME")
  dbPassword := os.Getenv("RDS_PASSWORD")

  // If no database credentials are provided, use in-memory storage
  if dbHost == "" || dbName == "" {
    fmt.Println("No database credentials found, using in-memory storage")
    dbInstance = &Store{
      db: nil,
      Q:  memory.New(),
    }
    fmt.Println("In-memory store initialized")
    return nil
  }

  dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
  InitDB(dataSourceName)
  fmt.Println(dataSourceName)
  return nil
}

func GetStore() (*Store, error) {
  if dbInstance == nil {
    return nil, fmt.Errorf("Store not initialized")
  }
	return dbInstance, nil
}
