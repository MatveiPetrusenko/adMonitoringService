package http

import (
	"encoding/json"
	"log"
	"net/http"

	"main.go/handler"
)

var RequestDataStream = make(chan requestData)

type requestData struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

func HandleRequests() {
	http.HandleFunc("/subscription", subscription)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func subscription(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var data requestData

	if request.Method == "POST" {
		writer.WriteHeader(http.StatusOK)

		err := json.NewDecoder(request.Body).Decode(&data)
		if err != nil {
			log.Fatalln("There was an error decoding the request body into the struct")
		}
	}

	errEmail := handler.CheckInputEmail(data.Email)
	if errEmail != nil {
		log.Fatalln(errEmail)
	}

	if !handler.CheckInputLink(data.Link) {
		return
	}

	RequestDataStream <- data
}
