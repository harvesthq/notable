package notable

import (
	"bytes"
	"fmt"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	"log"
	"strings"
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
	count := len(categoryNotes.Notes)
	announcements := pluralize(count, "Announcement")

	return fmt.Sprintf("#%s &mdash; %s", categoryNotes.Name, announcements)
}

func (note *Note) AbbreviatedAuthor() string {
	nameParts := strings.Fields(note.Author)

	if len(nameParts) > 0 {
		firstName := nameParts[0]
		lastName := nameParts[len(nameParts)-1]

		lastNameInitials := make([]string, 0)
		for _, name := range strings.Split(lastName, "-") {
			lastNameInitials = append(lastNameInitials, fmt.Sprintf("%s.", string(name[0])))
		}

		return fmt.Sprintf("%s %s", firstName, strings.Join(lastNameInitials, "-"))
	} else {
		return note.Author
	}
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
	subject := pluralize(len(Notes()), "Notable Announcement")

	message := &mandrill.Message{}
	message.AddRecipient("harvest.team@getharvest.com", "Harvest Team", "to")
	message.FromEmail = "notable@getharvest.com"
	message.FromName = "Harvest Notables"
	message.Subject = subject
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

func pluralize(count int, singularForm string) string {
	pluralForm := fmt.Sprintf("%s%s", singularForm, "s")

	if count == 1 {
		return fmt.Sprintf("1 %s", singularForm)
	} else {
		return fmt.Sprintf("%d %s", count, pluralForm)
	}
}
