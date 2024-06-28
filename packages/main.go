package main

import (
	"fmt"
	"math"
)

// func main() {
//	These inefered datatype
// 	var investmentAmount = 1000
// 	var execptedReturnRate = 5.5
// 	var years = 10

// 	var futureValue = float64(investmentAmount) * math.Pow(1+execptedReturnRate/100, float64(years))
// 	fmt.Println(futureValue)
// }

// func main() {
// 	const inflation = 5.2
// 	var investmentAmount float64 = 1000
// 	execptedReturnRate := 5.5 // := lets the go identify the data type
// 	var years float64 = 10

// 	// var investmentAmount, years float64 = 1000,10 => allowed
// 	// var investmentAmount float64, years float64 = 1000,10 => not allowed

// 	futureValue := investmentAmount * math.Pow(1+execptedReturnRate/100, years)
// 	futureRealValue := futureValue / math.Pow(1+inflation/100, years)

// 	fmt.Println("Future Value =>", futureValue)
// 	fmt.Println("Future Real Value =>", futureRealValue)
// }

const inflation = 5.2

func main() {
	var investmentAmount float64
	var execptedReturnRate float64
	var years float64

	fmt.Print("Enter Investment Amount: ")
	fmt.Scan(&investmentAmount) //this Scan is limited to single digit and words

	fmt.Print("Enter Expected Return: ")
	fmt.Scan(&execptedReturnRate) //& is pointer

	fmt.Print("Enter Years: ")
	fmt.Scan(&years) //& is pointer

	futureValue, futureRealValue := calculateFutureValue(investmentAmount, execptedReturnRate, years)

	// fmt.Println("Investment Amount =>", investmentAmount)
	fmt.Printf("Investment Amount => %v\n", investmentAmount)
	// fmt.Println("Future Value =>", math.Round(futureValue))
	fmt.Printf("Future Value => %.2f\n", futureValue)
	fmt.Println("Future Real Value =>", math.Round(futureRealValue))

	customFunc("Hello")
}

func calculateFutureValue(investmentAmount float64, execptedReturnRate float64, years float64) (float64, float64) {
	fv := investmentAmount * math.Pow(1+execptedReturnRate/100, years)
	rfv := fv / math.Pow(1+inflation/100, years)
	return fv, rfv
}
