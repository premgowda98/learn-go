package main

import (
	"log"

	"error-handling/errors"
)

func main() {
	err := stack3()
	if err != nil {
		log.Println("something went wrong", err)
		// errors.Print(err) // Use the improved Print function for structured error output
	}
}

func stack1() error {
	return errors.New("stack1 error")
}

func stack2() error {
	err := stack1()
	if err != nil {
		return errors.Wrap(err, "stack2 error")
	}
	return nil
}

func stack3() error {
	err := stack2()
	if err != nil {
		return errors.Wrap(err, "stack3 error")
	}
	return nil
}
