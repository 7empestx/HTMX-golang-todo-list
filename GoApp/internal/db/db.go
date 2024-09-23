package db

import (
  "fmt"
  "os"
  "database/sql"
  "sync"
  "github.com/7empestx/GoHTMXToDoList/internal/db/store/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
	Q  *storedb.Queries
}

var (
	dbInstance *Store
	once       sync.Once
	initErr    error
)

func Init() error {
  dbHost := os.Getenv("RDS_HOSTNAME")
  dbName := os.Getenv("RDS_DB_NAME")
  dbUser := os.Getenv("RDS_USERNAME")
  dbPassword := os.Getenv("RDS_PASSWORD")

  dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
  fmt.Println(dataSourceName)

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
			Q:  storedb.New(db),
		}

    fmt.Println("Database connection successful")
	})

	return initErr
}

func GetStore() (*Store, error) {
	if dbInstance == nil {
		return nil, fmt.Errorf("database not initialized, call InitDB first")
	}
	return dbInstance, nil
}
