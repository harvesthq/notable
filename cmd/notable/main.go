package main

import (
	"encoding/json"
	"fmt"
	notable "github.com/harvesthq/notable"
	"log"
	"net/http"
	"os"
)

type OKResponse struct {
	Text string `json:"text"`
}

type SummaryResponse struct {
	Notes []notable.Note `json:"notes"`
}

func getAndSetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	token := "13XIbBjtLeimhPIY36DWZvdR"
	incomingToken := request.Form.Get("token")

	if incomingToken == token {
		var response []byte
		var err error

		responseWriter.Header().Set("Content-Type", "application/json")

		if request.Method == "POST" {
			notable.Record(request.Form.Get("user_name"), request.Form.Get("trigger_word"), request.Form.Get("text"))
			response, err = json.Marshal(OKResponse{"I got this."})
		} else {
			response, err = json.Marshal(SummaryResponse{notable.Notes()})
		}

		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.Write(response)
	} else {
		fmt.Printf("Invalid token received: %s\n", incomingToken)
		http.Error(responseWriter, "Invalid token", http.StatusForbidden)
		return
	}
}

func clearHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		notable.Reset()
	}
}

func emailHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		notable.SendEmail(os.Getenv("API_KEY"))
	} else {
		responseWriter.Header().Set("Content-Type", "text/html")
		responseWriter.Write([]byte(notable.Email()))
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	notable.Reset()
	http.HandleFunc("/email", emailHandler)
	http.HandleFunc("/clear", clearHandler)
	http.HandleFunc("/", getAndSetHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
