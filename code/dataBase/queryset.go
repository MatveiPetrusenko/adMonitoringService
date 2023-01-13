package dataBase

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func CheckSaleID(id string) bool {
	db := ConnectDB()
	defer db.Close()

	var dataID string
	requestData := "SELECT ad_id FROM advertisement WHERE ad_id = $1"

	sqlStatementData := db.QueryRow(requestData, id)
	err := sqlStatementData.Scan(&dataID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
			return false
		} else {
			panic(err)
		}
	}

	return true
}

func CheckEmail(id, eMail string) bool {
	db := ConnectDB()
	defer db.Close()

	var dataEmail string

	sqlStatementData := db.QueryRow(fmt.Sprintf("SELECT email FROM \"%s\" WHERE email = $1", id), eMail)
	errH := sqlStatementData.Scan(&dataEmail)
	if errH != nil {
		if errH == sql.ErrNoRows {
			fmt.Println("Zero rows found")
			return false
		} else {
			panic(errH)
		}
	}

	return true
}

func AddSale(id, name, currency, price, link string) {
	db := ConnectDB()
	defer db.Close()

	requestData := "INSERT INTO advertisement (ad_id, name, currency, price, link) VALUES ($1,$2,$3,$4,$5)"

	if _, err := db.Exec(requestData, id, name, currency, price, link); err != nil {
		fmt.Println(err)
	}
}

func AddEmail(id, eMail string) {
	db := ConnectDB()
	defer db.Close()

	requestData := "INSERT INTO \"%s\" (email, timer) VALUES ($1,$2)"
	query := fmt.Sprintf(requestData, id)

	if _, err := db.Exec(query, eMail, "12h"); err != nil {
		fmt.Println(err)
	}
}

func CreateTableEmail(id string) {
	db := ConnectDB()
	defer db.Close()

	requestData := "CREATE TABLE \"%s\" (id SERIAL PRIMARY KEY, email VARCHAR(256), timer VARCHAR(1024));"
	query := fmt.Sprintf(requestData, id)

	if _, err := db.Exec(query); err != nil {
		fmt.Println(err)
	}
}

func SelectAdvertisement() *sql.Rows {
	db := ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT ad_id, name, currency, price, link  FROM advertisement;")
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func GetEmail(tableName string) []string {
	db := ConnectDB()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT email FROM \"%s\";", tableName))
	if err != nil {
		log.Fatal(err)
	}

	emails := make([]string, 0)

	for rows.Next() {
		var email string

		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}

		emails = append(emails, email)
	}

	return emails
}

func UpdateSale(id, name, currency, price, link string) {
	db := ConnectDB()
	defer db.Close()

	requestData := "UPDATE advertisement SET name = $1, currency = $2, price = $3, link = $4 WHERE ad_id = $5;"

	if _, err := db.Exec(requestData, name, currency, price, link, id); err != nil {
		fmt.Println(err)
	}
}
