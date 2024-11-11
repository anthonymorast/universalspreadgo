package universalspreadrule

import (
	"slices"
)

// TODO: these would probably be protobufs

type CFICode int

const (
	Put  CFICode = iota
	Call         = iota
)

type Side int

const (
	Buy  Side = iota
	Sell      = iota
)

type Instrument int

const (
	Option Instrument = iota
	Equity            = iota
)

type SymbolInfo struct {
	Underlier   string
	Symbol      string
	Price       float64
	Maturity    string
	Qty         float64
	StrikePrice float64
	TradeSide   Side
	PutCall     CFICode
	Instrument  Instrument
}

type marginRequest struct {
	Symbols          []SymbolInfo
	ShortMarginRate  float64
	LongMarginRate   float64
	OptionMarginRate float64
}

func NewDefaultMarginRequest() *marginRequest {
	return &marginRequest{[]SymbolInfo{}, 0.30, 0.25, 0.20}
}

func NewMarginRequest(longMarginRate, shortMarginRate, optionMarginRate float64) *marginRequest {
	return &marginRequest{[]SymbolInfo{}, shortMarginRate, longMarginRate, optionMarginRate}
}

func (request *marginRequest) AddOption(
	root string, symbol string, price float64, maturity string, qty float64, side Side, putCall CFICode, strike float64) {
	var info SymbolInfo = SymbolInfo{root, symbol, price, maturity, qty, strike, side, putCall, Option}
	request.Symbols = append(request.Symbols, info)
}

func (request *marginRequest) AddEquity(symbol string, price float64, qty float64, side Side) {
	var info SymbolInfo = SymbolInfo{symbol, symbol, price, "", qty, 0, side, Call, Equity}
	request.Symbols = append(request.Symbols, info)
}

func (request *marginRequest) GetEquityPosition(root string) *SymbolInfo {
	idx := slices.IndexFunc(request.Symbols, func(info SymbolInfo) bool { return info.Instrument == Equity && info.Symbol == root })

	if idx < 0 {
		return &SymbolInfo{}
	}
	return &request.Symbols[idx]
}

func (request *marginRequest) GetOptionPositions(root string) *[]SymbolInfo {
	var options []SymbolInfo
	for _, info := range request.Symbols {
		if info.Instrument == Option {
			options = append(options, info)
		}
	}
	return &options
}
