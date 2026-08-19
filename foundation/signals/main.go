package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 1)
	fmt.Println("Setting up signal notification...")
	signal.Notify(ch, os.Interrupt)
	// signal.Notify(ch, syscall.SIGINT)

	sig := <-ch
	fmt.Println("Received signal:", sig)
}
