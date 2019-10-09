package bitstamp

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const bitstampWsUrl = "wss://ws.bitstamp.net"
const wsTimeout = 60 * time.Second

type WsEvent struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type WsClient struct {
	ws       *websocket.Conn
	done     chan bool
	sendLock sync.Mutex
	Stream   chan *WsEvent
	Errors   chan error
}

func NewWsClient() (*WsClient, error) {
	c := WsClient{
		done:   make(chan bool, 1),
		Stream: make(chan *WsEvent),
		Errors: make(chan error),
	}

	// set up websocket
	ws, _, err := websocket.DefaultDialer.Dial(bitstampWsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket: %s", err)
	}
	c.ws = ws

	//
	// crux of the story
	//
	go func() {
		defer c.ws.Close()
		for {
			c.ws.SetReadDeadline(time.Now().Add(wsTimeout))
			select {
			case <-c.done:
				return
			default:
				var message []byte
				var err error
				_, message, err = c.ws.ReadMessage()
				if err != nil {
					c.Errors <- err
					continue
				}
				e := &WsEvent{}
				err = json.Unmarshal(message, e)
				if err != nil {
					c.Errors <- err
					continue
				}
				c.Stream <- e
			}
		}
	}()

	return &c, nil
}

func (c *WsClient) Close() {
	c.done <- true
}

func (c *WsClient) Subscribe(channels ...string) {
	for _, channel := range channels {
		sub := WsEvent{
			Event: "bts:subscribe",
			Data: map[string]interface{}{
				"channel": channel,
			},
		}
		c.sendEvent(sub)
	}
}

func (c *WsClient) Unsubscribe(channels ...string) {
	for _, channel := range channels {
		sub := WsEvent{
			Event: "bts:unsubscribe",
			Data: map[string]interface{}{
				"channel": channel,
			},
		}
		c.sendEvent(sub)
	}
}

func (c *WsClient) sendEvent(sub WsEvent) {
	c.sendLock.Lock()
	defer c.sendLock.Unlock()

	err := c.ws.WriteJSON(&sub)
	if err != nil {
		fmt.Println(err)
	}
}
