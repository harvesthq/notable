package notable

import (
	"bytes"
	"crypto/md5"
	"fmt"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	"io"
	"log"
	"strings"
	"text/template"
	"time"
)

type Note struct {
	Author    string    `json:"author"`
	Trigger   string    `json:"trigger"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

var notes []Note

func Record(author string, trigger string, text string) {
	notes = append(notes, Note{author, trigger, text, time.Now()})
}

func Notes() []Note {
	return notes
}

func Reset() {
	notes = make([]Note, 0)
}

func Email() string {
	var html bytes.Buffer
	funcMap := template.FuncMap{"gravatarHash": gravatarHash}

	notesTemplate, err := template.New("template.html").Funcs(funcMap).ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
	}

	err = notesTemplate.Execute(&html, notes)
	if err != nil {
		log.Fatal(err)
	}

	return html.String()
}

func EmailSubscribers(apiKey string) {
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

func gravatarHash(note Note) string {
	email := fmt.Sprintf("%s@getharvest.com", note.Author)
	hash := md5.New()
	io.WriteString(hash, strings.ToLower(email))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
