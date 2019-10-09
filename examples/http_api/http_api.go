package main

import (
	"bitstamp"
	"fmt"
	"log"
)

func main() {
	api := bitstamp.NewApiClient(
		bitstamp.Credentials("1", "ApiKey1", "api_key_secret"),
	)

	ticker, err := api.V1Ticker()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", ticker)

	ticker, err = api.V2Ticker("btcusd")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", ticker)

	ticker, err = api.V1HourlyTicker()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", ticker)

	ticker, err = api.V2HourlyTicker("btcusd")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", ticker)

	ob, err := api.V1OrderBook()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", ob)
}
