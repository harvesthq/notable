package main

import notable "github.com/harvesthq/notable"

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", notable.Handler)
	http.ListenAndServe(":"+port, nil)
}
