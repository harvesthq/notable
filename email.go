package notable

import (
	"bytes"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	"log"
	"text/template"
)

func Email() string {
	var html bytes.Buffer

	notesTemplate, err := template.ParseFiles("template.html")
	check(err)

	err = notesTemplate.Execute(&html, notesByCategory())
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
	grouped := make(map[string][]Note)

	for _, note := range Notes() {
		category = note.Category

		if len(grouped[category]) == 0 {
			grouped[category] = make([]Note, 1)
		}

		grouped[category] = append(grouped[category], note)
	}

	return grouped
}
