package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bitstonks/bitstamp-go/pkg/websocket"
)

func main() {
	c, err := websocket.NewWsClient()
	if err != nil {
		log.Panicf("error initializing client %v", err)
	}

	c.Subscribe("live_orders_btcusd", "live_trades_btcusd")

	go func() {
		for {
			select {
			case ev := <-c.Stream:
				fmt.Printf("%#v\n", ev)

			case err := <-c.Errors:
				fmt.Printf("--- ERROR: %#v\n", err)

			}
		}
	}()

	time.Sleep(3 * time.Second)

	fmt.Println("=== unsubscribing")
	c.Unsubscribe("live_orders_btcusd", "live_trades_btcusd")

	fmt.Println("=== sleeping")
	time.Sleep(1 * time.Second) // to clean up whatever...

	fmt.Println("=== closing")
	c.Close()
}
