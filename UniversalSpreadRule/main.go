package main

import (
	"fmt"
	"universalspreadrule/universalspreadrule"
)

func arrayPasser(arr any) {
	fmt.Println(arr)
}

func main() {
	var calculator universalspreadrule.ExcessMarginCalculator
	request := universalspreadrule.NewMarginRequest(0.25, 0.3, 0.20)

	// iron conndor
	// request.AddOption("INTC", "INTC 241227C30000", 0.75, "241220", 1, universalspreadrule.Buy, universalspreadrule.Call, 30.)
	// request.AddOption("INTC", "INTC 241220C25000", 1.25, "241220", 1, universalspreadrule.Sell, universalspreadrule.Call, 25.)
	// request.AddOption("INTC", "INTC 241220P22000", 1.0, "241220", 1, universalspreadrule.Sell, universalspreadrule.Put, 22.)
	// request.AddOption("INTC", "INTC 241220P10000", 0.90, "241220", 1, universalspreadrule.Buy, universalspreadrule.Put, 10.)

	request.AddOption("TEST", "TEST 241220C20000", 1.50, "241220", 1, universalspreadrule.Buy, universalspreadrule.Call, 90.)
	request.AddOption("TEST", "TEST 241227C30000", 0.75, "250117", 1, universalspreadrule.Buy, universalspreadrule.Call, 120.)
	request.AddOption("TEST", "TEST 241220C25000", 1.25, "250117", 1, universalspreadrule.Sell, universalspreadrule.Call, 100.)


	// request.AddEquity("INTC", 22.33, 100, universalspreadrule.Buy)

	results := calculator.CalculateOrderMargin(request)

	fmt.Println(results)
}
