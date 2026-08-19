package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string, m *sync.Mutex) {
	defer wg.Done()

	m.Lock()
	msg = s
	m.Unlock()
}

func printMessage() {
	fmt.Println(msg)
}

type Income struct {
	Source string
	Amount int
}

func main() {

	// msg = "Initial"

	// var mutex sync.Mutex

	// wg.Add(2)

	// go updateMessage("Hello Bengaluru", &mutex)

	// go updateMessage("Hello Mysuru", &mutex)

	// wg.Wait()

	// printMessage()

	var bankBalance int
	var mutex sync.Mutex

	fmt.Printf("Bank balance is %d", bankBalance)
	fmt.Println()

	incomes := []Income{
		{Source: "Job", Amount: 100},
		{Source: "Freelance", Amount: 300},
		{Source: "Partime", Amount: 50},
		{Source: "Invest", Amount: 100},
	}

	wg.Add(len(incomes))

	for _, inc := range incomes {
		go func(income Income) {
			defer wg.Done()
			for week := 1; week < 52; week++ {

				mutex.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				mutex.Unlock()

				fmt.Printf("On week %d, earned %d from %s", week, income.Amount, income.Source)
				fmt.Println()
			}
		}(inc)
	}

	wg.Wait()

	fmt.Printf("Final Balance %d", bankBalance)
	fmt.Println()
}
