package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gowda <message>")
		os.Exit(1)
	}

	// Print the message from command-line argument
	fmt.Println(os.Args[1])
}
