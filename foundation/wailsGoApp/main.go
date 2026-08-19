package main

import (
	"fmt"

	"github.com/premgowda/learn-go/wailsAppGo/ui"
)

func main() {
	fmt.Println("Hello, Wails!")

	fmt.Println("Starting UI from main.go")

	ui.StartUI()

	select {}
}
