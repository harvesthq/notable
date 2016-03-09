package main

import (
	notable "github.com/harvesthq/notable"
	"log"
	"os"
	"strconv"
)

func main() {
	if port, err := strconv.Atoi(os.Getenv("SMTP_PORT")); err == nil {
		if len(notable.Notes()) > 0 {
			notable.SendEmail(
				os.Getenv("SMTP_HOST"),
				port,
				os.Getenv("SMTP_USERNAME"),
				os.Getenv("SMTP_PASSWORD"),
				os.Getenv("FROM_EMAIL"),
				os.Getenv("FROM_NAME"),
				os.Getenv("TO_EMAIL"),
				os.Getenv("TO_NAME"),
			)
			if os.Getenv("NO_RESET") == "" {
				notable.Reset()
			} else {
				log.Print("Not resetting notes")
			}
		}
	} else {
		log.Fatal(err)
	}
}
