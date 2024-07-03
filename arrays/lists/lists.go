package lists

import "fmt"

type addArray struct {
	names  string
	arrays [4]string
}

func Array() {
	fmt.Println("Static Array")

	var arryDefined [4]string
	var arryDefinedint [4]int
	prices := [4]float64{10.5, 20.5, 3.6, 3.8} // specifying we want to store 4 float values

	arryDefinedint[2] = 2

	fmt.Println(arryDefinedint, len(arryDefinedint))
	fmt.Println(arryDefined, cap(arryDefined))
	fmt.Println(prices)

	//acessing values
	fmt.Println(prices[0:2], len(prices[0:2]), cap(prices[0:2]), cap(prices[1:3][:]))

	//Dynamic array
	fmt.Println("Dynamic Array")
	pricesDynamic := []float64{356} // dynamic array
	// otherArray := [2]float64{25.3, 26.5}

	newArray := append(pricesDynamic, 256) // append will not add to the same array but it will create new array

	fmt.Println(newArray)
	// fmt.Println(otherArray + otherArray) // not possible
	// fmt.Println(newArray + otherArray) // not possible

	arr := addArray{
		names:  "struct",
		arrays: [4]string{"kumar", "hello"},
	}

	arr.arrays[0] = "Prem"

	fmt.Println(arr)

	Slice_vs_array()
}

func Slice_vs_array() {
	//array
	fmt.Println("\n Array")
	x := [3]int{1, 2, 3}
	y := x

	fmt.Println(x)
	fmt.Println(y)
	y[2] = 4
	fmt.Println(x)
	fmt.Println(y)

	// slice
	fmt.Println("\n Slice")
	x1 := []int{2, 4}
	y1 := x1

	fmt.Println(x1)
	fmt.Println(y1)
	y1[0] = 4
	fmt.Println(x1)
	fmt.Println(y1)

}

func LoopArray() {
	values := [3]float64{1, 2, 3}

	for index, val := range values {
		fmt.Println(index, val)
	}
}
