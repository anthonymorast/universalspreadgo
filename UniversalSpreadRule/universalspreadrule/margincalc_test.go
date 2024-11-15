package universalspreadrule

import (
	"testing"
)

var calculator ExcessMarginCalculator

func TestLongPurchase(t *testing.T) {
	request := NewDefaultMarginRequest()
	request.AddEquity("TEST", 123.50, 125, Buy)

	results := calculator.CalculateOrderMargin(request)
	expected := 123.50 * .25 * 125
	if results.MarginReq != expected {
		t.Fatalf("expected req=%f\nactual=%f", expected, results.MarginReq)
	}
}

func TestCallCreditSpread(t *testing.T) {
	request := NewMarginRequest(0.25, 0.3, 0.20)

	// call credit spread strike difference = margin req = $500 - premium received = $500 - 125 + 75 = $450
	request.AddOption("INTC", "INTC 241227C30000", 0.75, "241220", 1, Buy, Call, 30.)
	request.AddOption("INTC", "INTC 241220C25000", 1.25, "241220", 1, Sell, Call, 25.)

	results := calculator.CalculateOrderMargin(request)

	if results.OptionPremium != -50. || results.OptionRequirement != 500. {
		t.Errorf("actual: margin req=%f, premium=%f, total=%f\nexpected margin req = 500 premium = -50 total = 450", results.OptionRequirement, results.OptionPremium, results.OptionRequirement+results.OptionPremium)
	}
}

func TestCallCreditSpreadWithLongOption(t *testing.T) {
	request := NewMarginRequest(0.25, 0.3, 0.20)

	request.AddOption("TEST", "TEST 241220C20000", 1.50, "241220", 1, Buy, Call, 90.)
	request.AddOption("TEST", "TEST 241227C30000", 0.75, "250117", 1, Buy, Call, 120.)
	request.AddOption("TEST", "TEST 241220C25000", 1.25, "250117", 1, Sell, Call, 100.)

	results := calculator.CalculateOrderMargin(request)

	if results.OptionPremium != 100. || results.OptionRequirement != 2000. {
		t.Errorf("actual: margin req=%f, premium=%f, total=%f\nexpected margin req = 500 premium = -50 total = 450", results.OptionRequirement, results.OptionPremium, results.OptionRequirement+results.OptionPremium)
	}
}

func TestIronCondor(t *testing.T) {
	request := NewDefaultMarginRequest()

	request.AddOption("INTC", "INTC 241227C30000", 0.75, "250117", 1, Buy, Call, 30.)
	request.AddOption("INTC", "INTC 241220C25000", 1.25, "250117", 1, Sell, Call, 25.)
	request.AddOption("INTC", "INTC 241220P22000", 1.0, "250117", 1, Sell, Put, 22.)
	request.AddOption("INTC", "INTC 241220P10000", 0.90, "250117", 1, Buy, Put, 10.)

	results := calculator.CalculateOrderMargin(request)

	if results.OptionPremium != -60. || results.OptionRequirement != 1200 {
		t.Errorf("actual: margin req=%f, premium=%f, total=%f\nexpected margin req = 1200 premium = -60 total = 1140", results.OptionRequirement, results.OptionPremium, results.OptionRequirement+results.OptionPremium)
	}
}
