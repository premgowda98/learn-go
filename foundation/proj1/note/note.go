package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Note struct {
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_at"`
}

func New(title, content string) (Note, error) {

	if title == "" || content == "" {
		return Note{}, errors.New("Invalid")
	}

	return Note{
		Title:       title,
		Content:     content,
		CreatedDate: time.Now(),
	}, nil
}

func (n Note) Display() {
	fmt.Println("Title: ", n.Title)
	fmt.Println("Content: ", n.Content)
}

func (n Note) Save() error {
	fileName := strings.ReplaceAll(n.Title, " ", "_")
	fileName = strings.ToLower(fileName) + ".json"

	json, err := json.Marshal(n) // this can write only those field which are exposed, i.e caps letter

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, json, 0644)
}
