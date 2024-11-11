package universalspreadrule

import (
	"math"
	"slices"
	"sort"
)

type MarginCalculator interface {
	CalculateOrderMargin(request marginRequest) CalculationResults
}

type ExcessMarginCalculator struct {
}

func (excessCalc ExcessMarginCalculator) walkCalls(callList []SymbolInfo, callMarginCh chan float64) {
	defer close(callMarginCh)

	if len(callList) < 2 { // single call or no call => no pair
		callMarginCh <- 0
		return
	}

	// sort by strike ascending
	sort.Slice(callList, func(idx1, idx2 int) bool {
		return callList[idx1].StrikePrice < callList[idx2].StrikePrice
	})

	callMargin := 0.
	lowIdx, highIdx := 0, 1
	lowStrikeInfo := &callList[lowIdx]
	highStrikeInfo := &callList[highIdx]

	for lowIdx < len(callList) && highIdx < len(callList) {
		if lowStrikeInfo.TradeSide == Sell { // short option
			// find higher strike long calls with GTE maturity
			if highStrikeInfo.Maturity >= lowStrikeInfo.Maturity && highStrikeInfo.TradeSide == Buy {
				// call credit spread
				qtyReduction := min(highStrikeInfo.Qty, lowStrikeInfo.Qty)
				highStrikeInfo.Qty -= qtyReduction
				lowStrikeInfo.Qty -= qtyReduction

				callMargin += lowStrikeInfo.StrikePrice - highStrikeInfo.StrikePrice
			}

		} else { // long option
			// find higher strike short calls with LTE maturity
			if highStrikeInfo.Maturity <= lowStrikeInfo.Maturity && highStrikeInfo.TradeSide == Sell {
				// reduce quantities, the trade is a debit spread, i.e. premium-only
				qtyReduction := min(highStrikeInfo.Qty, lowStrikeInfo.Qty)
				highStrikeInfo.Qty -= qtyReduction
				lowStrikeInfo.Qty -= qtyReduction
			}
		}

		if lowStrikeInfo.Qty == 0 {
			// get the next index with a pairable qty
			lowIdx = slices.IndexFunc(callList, func(info SymbolInfo) bool { return info.Qty > 0 })
			if lowIdx >= len(callList) || lowIdx < 0 {
				break
			}
			lowStrikeInfo = &callList[lowIdx]

			// start pairing at next available index
			highIdx = lowIdx + 1
			if highIdx >= len(callList) {
				break
			}
			highStrikeInfo = &callList[highIdx]
		} else {
			highIdx += 1
			if highIdx >= len(callList) {
				break
			}
			highStrikeInfo = &callList[highIdx]
		}
	}

	callMarginCh <- callMargin
}

func (excessCalc ExcessMarginCalculator) walkPuts(putList []SymbolInfo, putMarginCh chan float64) {
	defer close(putMarginCh)

	if len(putList) < 2 { // single option or no option => no pair
		putMarginCh <- 0
		return
	}

	// sort with strikes descending
	sort.Slice(putList, func(idx1, idx2 int) bool {
		return putList[idx1].StrikePrice > putList[idx2].StrikePrice
	})

	putMargin := 0.
	baseIdx, pairIdx := 0, 1
	baseInfo, pairInfo := &putList[baseIdx], &putList[pairIdx]
	for baseIdx < len(putList) && pairIdx < len(putList) {
		if baseInfo.TradeSide == Sell { // short put -> match with long put w/ > strike and LTE expiration
			if pairInfo.Maturity <= baseInfo.Maturity && pairInfo.TradeSide == Buy {
				pairQty := min(pairInfo.Qty, baseInfo.Qty)
				baseInfo.Qty -= pairQty
				pairInfo.Qty -= pairQty

				putMargin += pairInfo.StrikePrice - baseInfo.StrikePrice
			}
		} else { // long put
			pairQty := min(pairInfo.Qty, baseInfo.Qty)
			baseInfo.Qty -= pairQty
			pairInfo.Qty -= pairQty

			// putMargin += 0 // premium only debit spread
		}

		if baseInfo.Qty == 0 {
			// get the next option with a pairable qty
			baseIdx = slices.IndexFunc(putList, func(info SymbolInfo) bool { return info.Qty > 0 })
			pairIdx += 1

			if baseIdx >= len(putList) || pairIdx >= len(putList) || baseIdx < 0 {
				break
			}

			baseInfo = &putList[baseIdx]
			pairInfo = &putList[pairIdx]
		} else { // no pair, try next option
			pairIdx += 1
			if pairIdx >= len(putList) {
				break
			}
			pairInfo = &putList[pairIdx]
		}
	}

	putMarginCh <- putMargin
}

func (excessCalc ExcessMarginCalculator) universalSpreadRule(optionSymbols *[]SymbolInfo) float64 {
	putMarginCh := make(chan float64)
	callMarginCh := make(chan float64)

	callList, putList := []SymbolInfo{}, []SymbolInfo{}
	for _, info := range *optionSymbols {
		if info.PutCall == Call {
			callList = append(callList, info)
		} else {
			putList = append(putList, info)
		}
	}

	go excessCalc.walkCalls(callList, callMarginCh)
	go excessCalc.walkPuts(putList, putMarginCh)

	callMargin, putMargin := <-callMarginCh, <-putMarginCh

	// return min(<-putMarginCh, <-callMarginCh)
	return math.Abs(min(callMargin, putMargin)) * 100.
}

func (excessCalculator ExcessMarginCalculator) calculateNakedEquityMargin(longMarginRate, shortMarginRate, qty, price float64, side Side) float64 {
	if side == Buy { // long margin
		return longMarginRate * qty * price
	}

	if price < 5.0 {
		price = max(2.50, price)
		shortMarginRate = 1.0
		return price * qty * shortMarginRate * -1 // margin charge for negative qty
	}

	price = max(5.0, price)
	shortMarginRate = max(0.30, shortMarginRate)
	return price * qty * shortMarginRate * -1
}

func (excessCalculator ExcessMarginCalculator) getOrderPremium(request *marginRequest) float64 {
	var premium float64 = 0.
	for _, info := range request.Symbols {
		if info.Instrument == Option {
			var mult float64
			if mult = 1.0; info.TradeSide == Sell {
				mult = -1.0
			}
			premium += info.Price * info.Qty * 100. * mult
		}
	}
	return premium
}

func (excessCalc ExcessMarginCalculator) CalculateOrderMargin(request *marginRequest) CalculationResults {
	var results CalculationResults

	equityPosition := request.GetEquityPosition(request.Symbols[0].Underlier)
	unpairedEquityPosition := equityPosition.Qty
	optionPositions := request.GetOptionPositions(request.Symbols[0].Underlier)

	// calc order premium
	results.OptionPremium = excessCalc.getOrderPremium(request)

	// pair with equities and calc covered/married calls and puts
	// no pair? equity margin
	if equityPosition != nil {
		results.MarginReq += excessCalc.calculateNakedEquityMargin(request.LongMarginRate, request.ShortMarginRate, unpairedEquityPosition, equityPosition.Price, equityPosition.TradeSide)
	}

	// universal spread rule
	optionMargin := excessCalc.universalSpreadRule(optionPositions)

	// straddles/strangles

	// individual options

	results.OptionRequirement += optionMargin
	return results
}
