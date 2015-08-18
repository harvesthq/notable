package main

import (
	"encoding/json"
	"fmt"
	notable "github.com/harvesthq/notable"
	"log"
	"net/http"
	"net/url"
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
			recordNote(request.Form)
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
	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.Write([]byte(notable.Email()))
}

func recordNote(form url.Values) {
	user_id := form.Get("user_id")
	category := form.Get("trigger_word")
	text := form.Get("text")
	channel := form.Get("channel_name")
	slackToken := os.Getenv("SLACK_API_TOKEN")

	notable.Record(user_id, category, text, channel, slackToken)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/email", emailHandler)
	http.HandleFunc("/clear", clearHandler)
	http.HandleFunc("/", getAndSetHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
