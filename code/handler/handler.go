package handler

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

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

func ReadInputLink(link string) string {
	fmt.Scan(&link)
	return link
}

func CheckInputLink(link string) bool {
	return strings.Contains(link, "www.ebay.com/itm/")
}

func ReadInputEmail(eMail string) string {
	fmt.Scan(&eMail)
	return eMail
}

func CheckInputEmail(eMail string) bool {
	return strings.Contains(eMail, "@")
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
	// Host Email
	from := os.Getenv("HOSTEMAIL")
	password := os.Getenv("HOSTPASSWORD")

	// Receiver email address.
	to := dataBase.GetEmail(oldData.ID)

	// smtp server configuration.
	smtpHost := os.Getenv("HOST")
	smtpPort := os.Getenv("PORT")

	// Message.
	message := []byte("Subject: Notification\r\n" + "\r\n" + "Price" + "\r\n" + " Was:" + string(oldData.Price) + "\r\n" + "Become:" + string(newData.Price) + "\r\n")

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
