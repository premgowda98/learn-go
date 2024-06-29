package main

import (
	"bufio"
	"fmt"
	"learn/project1/note"
	"os"
	"strings"
)

func main() {

	title, content := getNoteData()

	userNote, err := note.New(title, content)

	if err != nil {
		fmt.Println(err)
		return
	}

	// userNote.Display()

	err = userNote.Save()

	if err != nil {
		fmt.Println("Saving failed")
		return
	}

	fmt.Println("Saved successfully")
}

func getNoteData() (string, string) {
	title := getUserInput("Note Title: ")

	content := getUserInput("Note Content: ")

	return title, content
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	// var value string
	// fmt.Scanln(&value) // Scan will not work for long text and also for space

	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n') // Since it is single character (rune) we should use '' not ""

	if err != nil {
		return ""
	}

	text = strings.TrimSuffix(text, "\n")

	return text
}
