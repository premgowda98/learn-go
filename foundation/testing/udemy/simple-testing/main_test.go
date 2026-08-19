package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
)

// manually testing all the scenarios
// func Test_isPrime(t *testing.T) {
// 	result, msg := isPrime(0)

// 	if result {
// 		t.Errorf("0 is not a prime number, expected false but recieved tru")
// 	}

// 	if msg != "0 is not prime by definitions"{
// 		t.Error("incorrect error message")
// 	}
// }

// Using table testns

func Test_isPrime(t *testing.T) {
	testTables := []struct {
		name     string
		input    int
		expected bool
		message  string
	}{
		{"prime", 7, true, "7 is a prime number"},
		{"not prime", 8, false, "8 is not prime"},
		{"not prime", 0, false, "0 is not prime by definition"},
		{"not prime", -5, false, "-5 is not prime since it is -ve"},
	}

	for _, test := range testTables {
		result, msg := isPrime(test.input)

		if test.expected && !result {
			t.Error("expected true but got false")
		}

		if !test.expected && result {
			t.Error("expected false but got true")
		}

		if msg != test.message {
			t.Error("unexpected message")
		}
	}
}

func Test_prompt(t *testing.T) {
	oldOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	prompt()

	_ = w.Close()

	os.Stdout = oldOut

	out, _ := io.ReadAll(r)

	if string(out) != "-> " {
		t.Error("incoorect prompt")
	}
}

func Test_intro(t *testing.T) {
	oldOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	intro()

	_ = w.Close()

	os.Stdout = oldOut

	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "Enter a whole number") {
		t.Error("incorect intro")
	}
}

func Test_checkNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1", "1 is not prime by definition"},
		{"q", ""},
		{"hi", "Please enter a proper number"},
	}

	for _, test := range tests {
		input := strings.NewReader(test.input)
		reader := bufio.NewScanner(input)

		res, _ := checkNumbers(reader)

		if !strings.EqualFold(res, test.expected) {
			t.Error("invalid output")
		}
	}
}
