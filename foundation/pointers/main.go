package main

import "fmt"

var age int = 32

func main() {
	fmt.Println("Age", age)
	fmt.Println("Age Pointer", &age)              // & providers memory address
	fmt.Println("Age Pointer destructure", *&age) // * provides value from address

	accessAge(&age)
	accessAge2(age)
	mutate(&age)
	fmt.Println("Age value", age)
}

func accessAge(age *int) int { // with *int we are specifying this function execpts pointer as input
	fmt.Println("Age pointer", age)
	return *age - 10 // here we are destructure the pointer to perform arthemetic operation
}

func accessAge2(age int) int {
	fmt.Println("Age var is copied?", &age)
	return age - 10
}

func mutate(age *int) {
	*age = *age + 10
	fmt.Println("Pointer after mutatation", age)
}
