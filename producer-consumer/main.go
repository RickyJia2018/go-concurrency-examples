package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const numOfPizzas = 10

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

	if pizzaNumber <= numOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false
		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++
		fmt.Printf("Making pizza number: %d. It will take %d seconds\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
		return &p
	}
	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}
func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0
	// run forever or until we receive a quit notification
	// try to make pizzas

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// try to make a pizza
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				// Close the channel
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
		// decision
	}
}

func main() {

	// seed random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// print out a message
	color.Cyan("The Pizzeria is open for businees!")
	color.Cyan("---------------------------------")
	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}
	// run the producer in background
	go pizzeria(pizzaJob)

	// create and run consumer
	for order := range pizzaJob.data {
		if order.pizzaNumber <= numOfPizzas {
			if order.success {
				color.Green(order.message)
				color.Green("Order #%d is out for delivery!", order.pizzaNumber)
			} else {
				color.Red(order.message)
				color.Red("Order #%d failed!", order.pizzaNumber)
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel")
			}
		}
	}
	// print out the ending message
	color.Cyan("--------------------------------")
	color.Cyan("Done for the day.")
	color.Cyan("Pizzas made: %d, but failed to make %d, with %d attempts in total", pizzasMade, pizzasFailed, total)
}
