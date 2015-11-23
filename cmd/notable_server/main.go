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
	token := request.Form.Get("token")

	if validToken(token) {
		var err error

		if request.Method == "POST" {
			recordNote(request.Form)
			responseWriter.Write([]byte("Got it!"))
		} else {
			err = respondWithJSON(responseWriter, SummaryResponse{notable.Notes()})
		}

		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Printf("Invalid token received: %s\n", token)
		http.Error(responseWriter, "Invalid token", http.StatusForbidden)
	}
}

func clearHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	token := request.Form.Get("token")

	if validToken(token) && request.Method == "POST" {
		notable.Reset()
	}
}

func emailHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	token := request.Form.Get("token")

	if validToken(token) {
		responseWriter.Header().Set("Content-Type", "text/html")
		responseWriter.Write([]byte(notable.Email()))
	}
}

func recordNote(form url.Values) {
	user_id := form.Get("user_id")
	category := form.Get("trigger_word")
	text := form.Get("text")
	slackToken := os.Getenv("SLACK_API_TOKEN")

	notable.Record(user_id, category, text, slackToken)
}

func slashCommandToken() string {
	return os.Getenv("SLACK_TOKEN")
}

func validToken(token string) bool {
	return token == slashCommandToken()
}

func respondWithJSON(responseWriter http.ResponseWriter, response interface{}) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(responseJSON)

	return nil
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
