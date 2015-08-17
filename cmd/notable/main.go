package main

import notable "github.com/harvesthq/notable"

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type OKResponse struct {
	Text string `json:"text"`
}

type SummaryResponse struct {
	Notes []string `json:"notes"`
}

func getAndSetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var response []byte
	var err error

	responseWriter.Header().Set("Content-Type", "application/json")

	if request.Method == "POST" {
		request.ParseForm()
		notable.Note(request.Form.Get("user_name"), request.Form.Get("text"))
		response, err = json.Marshal(OKResponse{"I got this."})
	} else {
		response, err = json.Marshal(SummaryResponse{notable.Summary()})
	}

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Write(response)
}

func clearHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		notable.Reset()
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	notable.Reset()
	http.HandleFunc("/notes", getAndSetHandler)
	http.HandleFunc("/notes/clear", clearHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
