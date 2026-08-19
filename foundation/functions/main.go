package main

import "fmt"

func main() {
	someVal := [4]int{1, 2, 3, 4}
	someVal2 := []int{1, 9, 3, 6, 2, 6}

	doubled := transForm(&someVal, douple)
	tripled := transForm(&someVal, triple)

	// anonyhmous function
	quadripled := transForm(&someVal, func(num int) int {
		return num * 4
	})

	// calling clouser function
	octiple := createTransformFunctions(8)
	octipled := transForm(&someVal, octiple)

	variadtic := variadicFunc(1, 9, 2, 3, 6)
	variadtic2 := variadicFunc(someVal2...)

	fmt.Println(doubled, tripled, quadripled, octipled, variadtic, variadtic2)
}

func transForm(value *[4]int, transformFunc func(int) int) []int {
	changedVal := []int{}

	for _, val := range value {
		changedVal = append(changedVal, transformFunc(val))
	}

	return changedVal
}

func douple(val int) int {
	return val * 2
}

func triple(val int) int {
	return val * 3
}

//Clouser functions

func createTransformFunctions(factor int) func(int) int {
	return func(number int) int {
		return number * factor
	}
}

//Variadic functions

func variadicFunc(num ...int) int {
	sum := 0

	for _, temp := range num {
		sum += temp
	}

	return sum
}
