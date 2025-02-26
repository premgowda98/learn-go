package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)

	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		fmt.Printf("Received Order #%d", pizzaNumber)
		fmt.Println()

		delayToMakePizza := rand.Intn(6) + 1
		fmt.Printf("Making pizza #%d, will be completed in %d seconds", pizzaNumber, delayToMakePizza)
		fmt.Println()

		time.Sleep(time.Duration(delayToMakePizza) * time.Second)

		rnd := rand.Intn(12) + 1
		msg := ""
		sucess := false

		if rnd < 5 {
			pizzasFailed++
			if rnd <= 2 {
				msg = fmt.Sprintf("Failed to make pizza #%d, ran out of ingredients", pizzaNumber)
			} else {
				msg = fmt.Sprintf("Failed to make pizza #%d, no power", pizzaNumber)
			}
		} else {
			pizzasMade++
			sucess = true
			msg = fmt.Sprintf("Pizza order #%d is ready", pizzaNumber)
		}

		total++

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     sucess,
		}

		return &p

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}

}

func pizzeria(pizzaMaker *Producer) {
	var i = 0

	for {
		currentPizza := makePizza(i)

		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <- pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	color.Cyan("Pizza Order is Online")
	color.Cyan("---------------------")

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(pizzaJob)

	for i := range pizzaJob.data{
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green("Order #%d out for delivery", i.pizzaNumber)
			} else {
				color.Red("Failed to make pizza #%d", i.pizzaNumber)
			}
		} else {
			color.Cyan("Done making pizza")

			err := pizzaJob.Close()

			if err != nil {
				color.Red("Failed to close channel", err)
			}
		}
	}

	color.Cyan("----------------")
	color.Cyan("Done for the day")
	color.Cyan("Made total of %d pizzas out of which %d got failed", total, pizzasFailed)
}
