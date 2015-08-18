package main

import (
	notable "github.com/harvesthq/notable"
	"os"
)

func main() {
	if len(notable.Notes()) > 0 {
		notable.SendEmail(os.Getenv("MANDRILL_API_KEY"))
		notable.Reset()
	}
}
