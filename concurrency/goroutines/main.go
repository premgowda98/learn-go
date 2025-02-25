package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}

var msg string

func updateMessage(s string) {
	msg = s
}

func printMessage() {
	fmt.Println(msg)
}

func main() {

	var wg sync.WaitGroup

	words := []string{"one", "two", "three", "four"}

	wg.Add(len(words))

	for _, word := range words {
		go printSomething(word, &wg)
	}

	wg.Wait()

	// go printSomething("with go-routine")
	wg.Add(1)
	printSomething("blocking", &wg)

	updateMessage("Hello Bengaluru")
	printMessage()

	updateMessage("Hello Mysuru")
	printMessage()
}
