package main

import (
	"fmt"
)

func main() {
	var input string

	fmt.Print("Enter the string: ")
	fmt.Scanf("%s\n", &input)

	var allwords []string
	numOfWords := 1

	tempIndex := 0

	for ind, char := range input {
		min, max := 'A', 'Z'
		if char >= min && char <= max {
			numOfWords++
			allwords = append(allwords, string(input[tempIndex:ind]))
			tempIndex = ind
		}
	}

	allwords = append(allwords, string(input[tempIndex:]))

	fmt.Printf("Total Number of words are: %d\nWords: %s\n", numOfWords, allwords)

}
