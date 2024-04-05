package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	if err := setup(); err != nil {
		log.Fatalf("setup failed: %v", err)
	} else {
		log.Println("setup finished successfully!")
	}
}
