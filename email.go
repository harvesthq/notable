package notable

import (
	"bytes"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	"log"
	"text/template"
	"time"
)

type Variables struct {
	Today           string
	NotesByCategory map[string][]Note
}

func Email() string {
	var html bytes.Buffer

	notesTemplate, err := template.ParseFiles("template.html")
	check(err)

	today := time.Now().Format("Monday January 2, 2006")
	variables := Variables{today, notesByCategory()}
	err = notesTemplate.Execute(&html, variables)
	check(err)

	return html.String()
}

func SendEmail(apiKey string) {
	client := mandrill.ClientWithKey(apiKey)

	message := &mandrill.Message{}
	message.AddRecipient("jason@getharvest.com", "Jason Dew", "to")
	message.AddRecipient("danny@getharvest.com", "Danny Wen", "to")
	message.FromEmail = "notable@getharvest.com"
	message.FromName = "Notable"
	message.Subject = "Daily Notable Digest"
	message.HTML = Email()

	_, err := client.MessagesSend(message)
	if err != nil {
		log.Print(err)
	}
}

func notesByCategory() map[string][]Note {
	var category string
	grouped := make(map[string][]Note, 0)

	for _, note := range Notes() {
		category = note.Category

		if len(grouped[category]) == 0 {
			grouped[category] = make([]Note, 0)
		}

		grouped[category] = append(grouped[category], note)
	}

	return grouped
}
