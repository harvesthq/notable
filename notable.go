package notable

import (
	"bytes"
	"encoding/json"
	"fmt"
	redis "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	mandrill "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/keighl/mandrill"
	slack "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/nlopes/slack"
	redisurl "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/soveran/redisurl"
	"log"
	"os"
	"strings"
	"text/template"
)

type Note struct {
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Trigger   string `json:"trigger"`
	Text      string `json:"text"`
}

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
	addNote(Note{authorName, avatarURL, trigger, text})
}

func Notes() []Note {
	db := redisConnection()
	defer db.Close()

	var notes []Note

	for _, id := range noteIDs(db) {
		notes = append(notes, fetchNote(db, id))
	}

	return notes
}

func Reset() {
	db := redisConnection()
	defer db.Close()

	for _, id := range noteIDs(db) {
		_, err := db.Do("DEL", noteKey(id))
		check(err)
	}

	_, err := db.Do("DEL", "notable:notes")
	check(err)
}

func Email() string {
	var html bytes.Buffer

	notesTemplate, err := template.ParseFiles("template.html")
	check(err)

	err = notesTemplate.Execute(&html, Notes())
	check(err)

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

func noteIDs(db redis.Conn) []int64 {
	var ids []int64
	rawIDs, err := redis.Values(db.Do("LRANGE", "notable:notes", 0, -1))
	check(err)

	redis.ScanSlice(rawIDs, &ids)

	return ids
}

func addNote(note Note) {
	db := redisConnection()
	defer db.Close()

	id, err := redis.Int64(db.Do("INCR", "notable:note"))
	check(err)

	_, err = db.Do("RPUSH", "notable:notes", id)
	check(err)

	_, err = db.Do("SET", noteKey(id), serialize(note))
	check(err)
}

func fetchNote(db redis.Conn, id int64) Note {
	noteAsJSON, err := redis.String(db.Do("GET", noteKey(id)))
	check(err)

	return deserialize(noteAsJSON)
}

func noteKey(id int64) string {
	return fmt.Sprintf("notable:note:%d", id)
}

func deserialize(noteAsJSON string) Note {
	var note Note

	err := json.Unmarshal([]byte(noteAsJSON), &note)
	check(err)

	return note
}

func serialize(note Note) string {
	noteAsJSON, err := json.Marshal(note)
	check(err)

	return string(noteAsJSON)
}

func redisConnection() redis.Conn {
	var connection redis.Conn
	var err error

	if len(os.Getenv("REDIS_URL")) > 0 {
		connection, err = redisurl.Connect()
	} else {
		connection, err = redis.Dial("tcp", ":6379")
	}

	check(err)

	return connection
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
