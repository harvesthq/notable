package notable

import (
	"bytes"
	"fmt"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	"log"
	"text/template"
	"time"
)

type Variables struct {
	Today           string
	NotesByCategory []CategoryNotes
}

type CategoryNotes struct {
	Name  string
	Notes []Note
}

func (categoryNotes *CategoryNotes) Title() string {
	size := len(categoryNotes.Notes)
	announcements := "Announcement"

	if size > 1 {
		announcements = announcements + "s"
	}

	return fmt.Sprintf("#%s &mdash; %d %s", categoryNotes.Name, size, announcements)
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
	message.AddRecipient("harvest.team@getharvest.com", "Harvest Team", "to")
	message.FromEmail = "notable@getharvest.com"
	message.FromName = "Notable"
	message.Subject = "Daily Notable Digest"
	message.HTML = Email()

	_, err := client.MessagesSend(message)
	if err != nil {
		log.Print(err)
	}
}

func notesByCategory() []CategoryNotes {
	var category string
	grouped := make(map[string]*CategoryNotes, 0)

	for _, note := range Notes() {
		category = note.Category

		if _, found := grouped[category]; !found {
			grouped[category] = &CategoryNotes{Name: category, Notes: make([]Note, 0)}
		}

		grouped[category].Notes = append(grouped[category].Notes, note)
	}

	categoryNotes := make([]CategoryNotes, 0)

	for _, value := range grouped {
		categoryNotes = append(categoryNotes, *value)
	}

	return categoryNotes
}
