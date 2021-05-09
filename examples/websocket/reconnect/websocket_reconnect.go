package main

import (
	"fmt"
	"log"
	"time"

	"github.com/samotarnik/bitstamp-go"
)

// Following app is an example of handling reconnect request from Websocket server. Note, that
// this shows the semantics but
type App struct {
	client           *bitstamp.WsClient
	requestNewClient chan struct{}
	close            chan struct{}
}

// A very simple example of handling reconnect event. Every reconnect event is handled in such a way, that can
// ensure no data loss.
func NewApp() App {
	c, err := bitstamp.NewWsClient()
	if err != nil {
		log.Panicf("error initializing client %v", err)
	}
	return App{
		client:           c,
		requestNewClient: make(chan struct{}),
		close:            make(chan struct{}),
	}
}

func (a *App) Run() {
	a.client.Subscribe("live_orders_btcusd", "live_trades_btcusd")
	for {
		select {
		case ev := <-a.client.Stream:
			fmt.Printf("%#v\n", ev)
			if a.client.IsReconnectRequest(ev) {
				a.requestNewClient <- struct{}{}
			}

		case err := <-a.client.Errors:
			fmt.Printf("--- ERROR: %#v\n", err)
		case <-a.close:
			// Wait before killing the app, as new connection must first be established
			time.Sleep(2 * time.Second)
			return
		}
	}
}

func main() {
	for {
		app := NewApp()
		go app.Run()
		<-app.requestNewClient
		app.close <- struct{}{}
	}
}
