package main

import notable "github.com/harvesthq/notable"

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", notable.Handler)
	http.ListenAndServe(":5000", nil)
}
