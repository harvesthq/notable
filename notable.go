package notable

import (
	"fmt"
	slack "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/nlopes/slack"
	"log"
	"regexp"
	"strings"
	"os"
)


func targetRoom() string {
  if room := os.Getenv("SLACK_CHANNEL"); len(room) > 0 {
    return room
  }

  return "general"
}

func Record(authorID string, category string, text string, slackToken string) {
	var authorName, avatarURL string

	api := slack.New(slackToken)
	user, err := api.GetUserInfo(authorID)
	text, category = effectiveTextAndCategory(text, category)

	if err == nil {
		authorName = user.Profile.RealName
		avatarURL = user.Profile.Image48
	} else {
		fmt.Printf("Error getting author information from Slack: %s\n", err)
		authorName = authorID
		avatarURL = ""
	}

	text = strings.TrimSpace(strings.TrimPrefix(text, category))
	note := Note{authorName, avatarURL, category, text}

	AddNote(note)
	notifyRoom(api, note)
}

func effectiveTextAndCategory(text string, category string) (string, string) {
	hashtag := extractHashtag(text)

	if len(hashtag) > 0 {
		category = hashtag
		text = strings.TrimSuffix(text, hashtag)
	}

	if len(category) == 0 {
		category = "notable"
	}

	return text, strings.Trim(category, ":# ")
}

func extractHashtag(text string) string {
	pattern := regexp.MustCompile(` #\w+\z`)

	return pattern.FindString(text)
}

func notifyRoom(api *slack.Client, note Note) {
	var attachments []slack.Attachment
	avatar := slack.Attachment{Text: note.Text, Color: "good"}
	heading := fmt.Sprintf("%s Announcement from %s", strings.Title(note.Category), note.Author)
	params := slack.PostMessageParameters{Username: heading, IconURL: note.AvatarURL, Attachments: append(attachments, avatar)}
	_, _, err := api.PostMessage(fmt.Sprintf("#%s", targetRoom()), "", params)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
