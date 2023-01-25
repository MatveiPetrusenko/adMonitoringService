// Package dataBase represent connecting to Postgresql Database /*
package dataBase

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDB represent connecting to Postgresql Database
func ConnectDB() *sql.DB {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgresdb", 5432, "postgres", "postgres", "postgres")

	db, errdb := sql.Open("postgres", postgresInfo)
	if errdb != nil {
		panic(errdb)
	}

	return db
}
