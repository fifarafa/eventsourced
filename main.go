package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
}

func openMySQLConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func createStream(db *sql.DB) error {
	panic("not implemented")
}
