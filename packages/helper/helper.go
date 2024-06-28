package helper

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

func SomeHelp() {
	fmt.Println("From other package")
	fmt.Println(randomdata.IpV4Address())
}
