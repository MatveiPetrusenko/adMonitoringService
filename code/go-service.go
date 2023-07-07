package main

import (
	"sync"

	"main.go/http"

	"main.go/scraper"

	"main.go/handler"
)

func main() {
	var sync sync.WaitGroup
	sync.Add(4)

	//
	go func() {
		defer sync.Done()

		for {
			http.HandleRequests()
		}
	}()

	//
	go func() {
		defer sync.Done()

		s := handler.NewScraper()
		for {
			requestData := <-http.RequestDataStream
			if requestData.Link == "" && requestData.Email == "" {
				continue
			}

			data, _ := s.Visit(requestData.Link)
			// todo: handling error

			handler.CheckExist(data, requestData.Email)
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
