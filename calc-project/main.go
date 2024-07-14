package main

import "project/calculator/prices"

func main() {
	taxrates := []float64{.26, .36}

	for _, taxrate := range taxrates {
		obj := prices.NewTaxRate(taxrate)
		obj.LoadData()
		obj.Process()
	}

}
