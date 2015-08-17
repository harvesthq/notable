package main

import (
	notable "github.com/harvesthq/notable"
	"os"
)

func main() {
	notable.SendEmail(os.Getenv("MANDRILL_API_KEY"))
}
