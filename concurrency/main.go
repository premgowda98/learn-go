package main

import (
	"fmt"
	"time"
)

func greet(pharse string) {
	fmt.Println("Hello", pharse)
}

func slowGreet(pharse string, doneChan chan bool) {
	time.Sleep(3 * time.Second)
	fmt.Println("Hello", pharse)

	doneChan <- true
}

// func main() {
// 	greet("Nice to meet you")
// 	greet("How are you")
// 	slowGreet("How are you................")
// 	greet("I am good")
// }

// adding go runs the fnction run in parallel
// but nothing is printed in the console
// snce all are running in parallel and the main function is done executing
// func main() {
// 	fmt.Println("Start")

// 	go greet("Nice to meet you") // adding go will run in parallel
// 	go greet("How are you")
// 	go slowGreet("How are you................")

// 	fmt.Println("Middle")

// 	go greet("I am good")

// 	fmt.Println("End")
// }

// channels

func main() {
	// go greet("Nice to meet you")
	// go greet("How are you")
	done := make(chan bool)
	go slowGreet("How are you................", done)
	fmt.Println(<-done)
	// go greet("I am good")
}
