package db

import (
  "fmt"
  "os"
  "database/sql"
  "github.com/7empestx/GoHTMXToDoList/internal/db/store/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
	Q  *storedb.Queries
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
		Q:  storedb.New(db),
	}

	return nil
}

func Init() error {
  dbHost := os.Getenv("RDS_HOSTNAME")
  dbName := os.Getenv("RDS_DB_NAME")
  dbUser := os.Getenv("RDS_USERNAME")
  dbPassword := os.Getenv("RDS_PASSWORD")

  dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
  InitDB(dataSourceName)
  fmt.Println(dataSourceName)
  return nil
}

func GetStore() (*Store, error) {
	return dbInstance, nil
}
