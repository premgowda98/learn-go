package main

import (
	"fmt"
	"io"
	"os"
)

func main(){
	data, err := readCSVFile("./problems.csv")

	if err !=nil{
		fmt.Println("Failed to read data")
	}

	fmt.Println(data)
	
}

func readCSVFile(filename string)([]byte, error){
	f, err := os.Open(filename)

	if err !=nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(f)

	if err !=nil {
		return nil, err
	}

	return data, nil

}