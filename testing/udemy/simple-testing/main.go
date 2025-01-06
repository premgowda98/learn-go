package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	intro()

	doneChanel := make(chan bool)

	go readUserInput(doneChanel)

	<-doneChanel

	fmt.Println("Thank you")

}

func readUserInput(doneChanel chan bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChanel <- true
			return
		}

		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	scanner.Scan()

	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	numToCheck, err := strconv.Atoi(scanner.Text())

	if err != nil {
		return "Please enter a proper number", false
	}

	_, msg := isPrime(numToCheck)

	return msg, false

}

func intro() {
	fmt.Println("Is it a prime number?")
	fmt.Println("---------------------")
	fmt.Println("Enter a whole number to check if its a whole number")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func isPrime(num int) (bool, string) {
	if num == 0 || num == 1 {
		return false, fmt.Sprintf("%d is not prime by definition", num)
	}

	if num < 0 {
		return false, fmt.Sprintf("%d is not prime since it is -ve", num)
	}

	for i := 2; i <= num/2; i++ {
		if num%i == 0 {
			return false, fmt.Sprintf("%d is not prime", num)
		}
	}

	return true, fmt.Sprintf("%d is a prime number", num)
}
