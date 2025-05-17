package main

import (
	"fmt"
	"runtime"
	_ "go.uber.org/automaxprocs"
)

func main() {
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(-1))
	fmt.Println("Max CPU: ", runtime.NumCPU())

}
