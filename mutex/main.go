package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {

	var bankBalance int
	var balance sync.Mutex
	fmt.Printf("initial account balance: %d\n", bankBalance)

	incomes := []Income{
		{
			Source: "Main Job",
			Amount: 500,
		},
		{
			Source: "Gift",
			Amount: 10,
		},
		{
			Source: "Partime Job",
			Amount: 50,
		},
		{
			Source: "Investment",
			Amount: 100,
		},
	}
	wg.Add(len(incomes))

	for i, income := range incomes {

		go func(i int, income Income) {
			defer wg.Done()
			for week := 1; week < 52; week++ {
				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				fmt.Printf("On week %d, you earned $%d from %s\n", i, income.Amount, income.Source)
				balance.Unlock()
			}
		}(i, income)
	}
	wg.Wait()

	fmt.Printf("final bank balance %d\n", bankBalance)
}
