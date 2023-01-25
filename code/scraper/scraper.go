package scraper

import (
	"fmt"
	"path"
	"strings"
	"unicode"

	"github.com/gocolly/colly"
)

type Data struct {
	ID       string
	Name     string
	Currency string
	Price    string
	Link     string
}

type Scrapper struct {
	collector  *colly.Collector
	data       *Data
	parseError error
}

// parseUrl getting advertising ID from link
func (data *Data) parseUrl(link string) {
	base := path.Base(link)

	for _, val := range base {
		if unicode.IsDigit(val) != true {
			break
		}

		data.ID += string(val)
	}
}

func (data *Data) setPrice(price string) error {
	stringData := strings.SplitN(price, " ", 2)
	// TODO: error handling

	data.Currency = stringData[0]
	data.Price = stringData[1]
	return nil
}

func (data *Data) addLink(link string) {
	data.Link = link
}

func (data *Data) setName(name string) {
	data.Name = strings.Trim(name, "—")
}

func NewScrapper() *Scrapper {
	c := colly.NewCollector()
	sc := &Scrapper{
		collector: c,
		data:      &Data{},
	}
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Accept-Language", "en-US")
		fmt.Println("-----------------------------")
		fmt.Println("Scraping:", request.URL)

		for key, value := range *request.Headers {
			fmt.Printf("%s: %s\n", key, value)
		}

		fmt.Println(request.Method)
	})

	c.OnHTML("div#LeftSummaryPanel", func(htmlElement *colly.HTMLElement) {
		price := htmlElement.ChildText("div.x-price-primary")
		// todo Price not empty, if not => set parseErr
		name := htmlElement.ChildText("h1.x-item-title__mainTitle")
		sc.data.setName(name)
		sc.data.setPrice(price)
	})

	c.OnResponse(func(response *colly.Response) {
		fmt.Println("-----------------------------")
		fmt.Println("Status:", response.StatusCode)
	})

	return sc
}

func (s *Scrapper) Visit(link string) (Data, error) {
	err := s.collector.Visit(link)
	if err != nil {
		return Data{}, err
	}

	if s.parseError != nil {
		// todo handing errors
	}
	s.data.addLink(link)
	s.data.parseUrl(link)

	return *s.data, nil
}
