package main

import (
	"fmt"
	"github.com/bitstonks/bitstamp-go/pkg/http"
	"log"
)

func main() {
	api := http.NewHttpClient(
		http.Credentials("1", "invalid", "invalid"),
	)

	// public endpoints
	ticker1, err := api.V2Ticker("btcusd")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("TICKER: %+v\n", ticker1)

	ticker2, err := api.V2HourlyTicker("btcusd")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("HOURLY TICKER: %+v\n", ticker2)

	ob, err := api.V2OrderBook("btcusd", 2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("ORDER BOOK - HIGHEST BID: %+v\n", ob.Bids[0])
	fmt.Printf("ORDER BOOK - LOWEST ASK: %+v\n", ob.Asks[0])

	txs, err := api.V2Transactions("btcusd", "hour")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("TRANSACTIONS: %+v\n", txs[0])

	info, err := api.V2TradingPairsInfo()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("TRADING PAIRS: %+v\n", info[0])

	ohlc, err := api.V2Ohlc("btcusd", 60, 2, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("CANDLES: %+v\n", ohlc.Data.Candles)

	eurusd, err := api.V2EurUsd()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("EURUSD: %+v\n", eurusd)
}
