package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := openMySQLConnection()
	if err != nil {
		panic(err)
	}
	err = createStream(db)
	if err != nil {
		panic(err)
	}
	log.Println("finished successfully!")
}

func openMySQLConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "local:local@tcp(127.0.0.1:3306)/local")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func createStream(db *sql.DB) error {
	_, err := db.Exec(`
  		CREATE TABLE IF NOT EXISTS stream(
  		    id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
  		    type VARCHAR(255) NOT NULL,
  		    version BIGINT NOT NULL
  		)`,
	)
	if err != nil {
		return fmt.Errorf("failed to create stream table: %w", err)
	}
	return nil
}
