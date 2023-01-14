package main

import (
	"fmt"
	"sync"

	"main.go/scraper"

	"main.go/handler"
)

func main() {
	var sync sync.WaitGroup
	sync.Add(3)

	//
	go func() {
		defer sync.Done()

		s := handler.NewScraper()
		for {
			var link, eMail string

			//
			link = handler.ReadInputLink(link)
			resultLink := handler.CheckInputLink(link)
			if !resultLink {
				fmt.Println("Incorrect input link")
				continue
			}

			//
			eMail = handler.ReadInputEmail(eMail)
			resultEmail := handler.CheckInputEmail(eMail)
			if !resultEmail {
				fmt.Println("Incorrect input Email")
				continue
			}

			data, _ := s.Visit(link)
			// todo: handling error

			handler.CheckExist(data, eMail)
		}
	}()

	//
	oldDataStream := make(chan handler.OldData)
	newDataStream := make(chan scraper.Data)

	//
	go func() {
		defer sync.Done()

		for {
			oldData, newData, result := handler.MonitorChanges()
			if result != false {
				oldDataStream <- oldData
				newDataStream <- newData
			}
		}
	}()

	//
	go func() {
		defer sync.Done()

		for {
			oldData := <-oldDataStream
			newData := <-newDataStream

			smtpHost, smtpPort, auth, from, to, message := handler.PrepareMessage(oldData, newData)
			handler.SendMessage(smtpHost, smtpPort, auth, from, to, message)

			handler.UpdateAdvertisement(newData)
		}
	}()

	sync.Wait()
}
