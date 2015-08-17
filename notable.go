package notable

import (
	"fmt"
)

var notes []string

func Note(author string, note string) {
	fmt.Printf("author=%s note=%s\n", author, note)
	notes = append(notes, note)
}

func Summary() []string {
	return notes
}

func Reset() {
	notes = make([]string, 0)
}
