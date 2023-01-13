package handler

import (
	"fmt"
	"log"
	"net/smtp"

	"main.go/dataBase"
	"main.go/scraper"
)

type OldData struct {
	ID       string
	Name     string
	Currency string
	Price    string
	Link     string
}

func NewScraper() *scraper.Scrapper {
	return scraper.NewScrapper()
}

func CheckExist(data scraper.Data, eMail string) {
	resultID := dataBase.CheckSaleID(data.ID)
	if !resultID {
		id := data.ID
		name := data.Name
		currency := data.Currency
		price := data.Price
		link := data.Link

		dataBase.AddSale(id, name, currency, price, link)
		dataBase.CreateTableEmail(id)
		dataBase.AddEmail(id, eMail)

		fmt.Println("Add sale in table")
	} else {
		fmt.Println("Sale already exist")
	}

	resultEmail := dataBase.CheckEmail(data.ID, eMail)
	if !resultEmail {
		id := data.ID

		dataBase.AddEmail(id, eMail)
		fmt.Println("Add eMail in table")
	} else {
		fmt.Println("Email already exist")
	}
}

func MonitorChanges() (OldData, scraper.Data, bool) {
	rows := dataBase.SelectAdvertisement()
	s := scraper.NewScrapper()

	var oldData OldData

	var result bool

	for rows.Next() {
		if err := rows.Scan(&oldData.ID, &oldData.Name, &oldData.Currency, &oldData.Price, &oldData.Link); err != nil {
			log.Fatal(err)
		}

		newData, _ := s.Visit(oldData.Link)

		if newData.ID == oldData.ID && newData.Name == oldData.Name && newData.Currency == oldData.Currency && newData.Price == oldData.Price && newData.Link == oldData.Link {
			fmt.Println("YES")
			continue
		} else {
			result = true
			return oldData, newData, result
		}
	}

	return oldData, scraper.Data{}, result
}

func PrepareMessage(oldData OldData, newData scraper.Data) (string, string, smtp.Auth, string, []string, []byte) {
	// Sender data
	// General Email
	from := "from@gmail.com"
	password := "<Email Password>"

	// Receiver email address.
	to := dataBase.GetEmail(oldData.ID)

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtpHost, smtpPort, auth, from, to, message

}

func SendMessage(smtpHost, smtpPort string, auth smtp.Auth, from string, to []string, message []byte) {
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func UpdateAdvertisement(newData scraper.Data) {
	dataBase.UpdateSale(newData.ID, newData.Name, newData.Currency, newData.Price, newData.Link)
}
