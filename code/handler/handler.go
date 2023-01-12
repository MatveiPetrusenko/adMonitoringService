package handler

import (
	"fmt"
	"log"

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

//newData.ID == OldData.ID && newData.Name == OldData.Name && newData.Currency == currency && newData.Price == price && newData.Link == link
//var adID, name, currency, price, link string
