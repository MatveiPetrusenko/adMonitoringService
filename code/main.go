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
			fmt.Scan(&link, &eMail)

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

		}
	}()

	sync.Wait()
}
