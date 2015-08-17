package notable

import (
	"bytes"
	"fmt"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	slack "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/nlopes/slack"
	"log"
	"strings"
	"text/template"
	"time"
)

type Note struct {
	Author    string    `json:"author"`
	AvatarURL string    `json:"avatar_url"`
	Trigger   string    `json:"trigger"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

var notes []Note

func Record(authorID string, trigger string, text string, slackToken string) {
	var authorName, avatarURL string

	api := slack.New(slackToken)
	user, err := api.GetUserInfo(authorID)

	if err == nil {
		authorName = user.Profile.RealName
		avatarURL = user.Profile.Image24
	} else {
		fmt.Printf("Error getting author information from Slack: %s\n", err)
		authorName = authorID
		avatarURL = ""
	}

	text = strings.TrimSpace(strings.TrimPrefix(text, trigger))
	notes = append(notes, Note{authorName, avatarURL, trigger, text, time.Now()})
}

func Notes() []Note {
	return notes
}

func Reset() {
	notes = make([]Note, 0)
}

func Email() string {
	var html bytes.Buffer

	notesTemplate, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
	}

	err = notesTemplate.Execute(&html, notes)
	if err != nil {
		log.Fatal(err)
	}

	return html.String()
}

func SendEmail(apiKey string) {
	client := mandrill.ClientWithKey(apiKey)

	message := &mandrill.Message{}
	message.AddRecipient("jason@getharvest.com", "Jason Dew", "to")
	message.FromEmail = "notable@getharvest.com"
	message.FromName = "Notable"
	message.Subject = "Daily Notable Digest"
	message.HTML = Email()

	_, err := client.MessagesSend(message)
	if err != nil {
		log.Print(err)
	}
}
