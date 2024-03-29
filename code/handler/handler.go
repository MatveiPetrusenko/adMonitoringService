package handler

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"
	"strings"

	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"

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

func CheckInputLink(link string) error {
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
		return err
	}

	urlWithoutPath := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	urlValidation := validation.NewStringRule(func(value string) bool {
		return strings.HasPrefix(link, value)
	}, "must start with the specified URL")

	err = validation.Validate(urlWithoutPath, urlValidation)
	if err != nil {
		fmt.Println("URL validation error:", err)
		return err
	} else {
		fmt.Println("URL is valid")
		return nil
	}

	//return strings.Contains(link, "www.ebay.com/itm/")
}

func CheckInputEmail(eMail string) error {
	return validation.Validate(&eMail, validation.Required, is.Email)
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
