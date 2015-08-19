package notable

import (
	"fmt"
	slack "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/nlopes/slack"
	"log"
	"regexp"
	"strings"
)

func Record(authorID string, category string, text string, channel string, slackToken string) {
	var authorName, avatarURL string

	api := slack.New(slackToken)
	user, err := api.GetUserInfo(authorID)
	text, category = effectiveTextAndCategory(text, category)

	if err == nil {
		authorName = user.Profile.RealName
		avatarURL = user.Profile.Image24
	} else {
		fmt.Printf("Error getting author information from Slack: %s\n", err)
		authorName = authorID
		avatarURL = ""
	}

	text = strings.TrimSpace(strings.TrimPrefix(text, category))
	note := Note{authorName, avatarURL, category, text}

	AddNote(note)

	if channel != "testbotroom" {
		notifyRoom(api, note)
	}
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
	avatar := slack.Attachment{AuthorIcon: note.AvatarURL}
	var attachments []slack.Attachment
	params := slack.PostMessageParameters{Username: note.Author, Attachments: append(attachments, avatar)}
	_, _, err := api.PostMessage("#testbotroom", note.Text, params)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
