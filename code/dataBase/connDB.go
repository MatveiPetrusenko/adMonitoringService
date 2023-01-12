// Package dataBase represent connecting to Postgresql Database /*
package dataBase

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var connStr = "user=postgres password=postgres dbname=postgres sslmode=disable"

// ConnectDB represent connecting to Postgresql Database
func ConnectDB() *sql.DB {
	db, errdb := sql.Open("postgres", connStr)
	if errdb != nil {
		panic(errdb)
	}

	return db
}
