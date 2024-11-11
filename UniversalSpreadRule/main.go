package main

import (
	"fmt"
	"universalspreadrule/universalspreadrule"
)

func arrayPasser(arr any) {
	fmt.Println(arr)
}

func main() {

	// input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	// output := make(chan int)
	// go func() {
	// 	for _, val := range input {
	// 		output <- val
	// 	}
	// 	close(output)
	// }()

	var calculator universalspreadrule.ExcessMarginCalculator
	request := universalspreadrule.NewMarginRequest(0.25, 0.3, 0.20)

	int_arr := []int{1, 2, 3, 4}
	fmt.Println(int_arr)
	arrayPasser(int_arr)

	// iron connie
	request.AddOption("INTC", "INTC 241227C30000", 0.75, "241220", 1, universalspreadrule.Buy, universalspreadrule.Call, 30.)
	request.AddOption("INTC", "INTC 241220C25000", 1.25, "241220", 1, universalspreadrule.Sell, universalspreadrule.Call, 25.)
	request.AddOption("INTC", "INTC 241220P22000", 1.0, "241220", 1, universalspreadrule.Sell, universalspreadrule.Put, 22.)
	request.AddOption("INTC", "INTC 241220P10000", 0.90, "241220", 1, universalspreadrule.Buy, universalspreadrule.Put, 10.)

	// request.AddEquity("INTC", 22.33, 100, universalspreadrule.Buy)

	results := calculator.CalculateOrderMargin(request)

	fmt.Println(results)

	// for range input {
	// 	fmt.Print(<-output, " ")
	// }
}
