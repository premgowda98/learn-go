package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	data, err := readCSVFile("problems.csv")

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	reader, err := parseCSV(data)
	if err != nil {
		fmt.Println("Error creating CSV reader:", err)
		return
	}
	records := processCSV(reader)

	fmt.Println("Math Quiz")

	correctAns := 0

	for ind, record := range records {
		fmt.Println(ind+1, "Ques:", record[0])
		userInput := takeUserInput()

		if userInput == record[1] {
			correctAns += 1
		}
	}

	fmt.Printf("\nYou got %d correct answers out of %d\n", correctAns, len(records))

}

func readCSVFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(f)

	if err != nil {
		return nil, err
	}

	return data, nil

}

func parseCSV(data []byte) (*csv.Reader, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	return reader, nil
}

func processCSV(reader *csv.Reader) [][]string {
	record, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading CSV data:", err)
	}

	return record
}

func takeUserInput() string {
	var userInput string

	fmt.Print("Ans: ")
	fmt.Scanf("%s", &userInput)

	return userInput
}
