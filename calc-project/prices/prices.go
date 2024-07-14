package prices

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

type TaxRatePrice struct {
	Prices     []float64
	TaxRate    float64
	TaxApplied map[string]float64
}

func (job *TaxRatePrice) LoadData() {
	file, err := os.Open("prices.txt")

	if err != nil {
		fmt.Println("Cloud not open file")
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(file)

	var lines []float64

	for scanner.Scan() {
		converted, _ := strconv.ParseFloat(scanner.Text(), 64)
		lines = append(lines, converted)
	}

	job.Prices = lines
	file.Close()
}

func (job *TaxRatePrice) Process() {
	job.LoadData()

	results := make(map[string]float64)

	for _, price := range job.Prices {
		results[fmt.Sprintf("%.2f", price)] = math.Round(price * (1 + job.TaxRate))
	}

	fmt.Println(results)
}

func NewTaxRate(taxrate float64) TaxRatePrice {
	return TaxRatePrice{
		Prices:  []float64{100, 12, 45},
		TaxRate: taxrate,
	}
}
