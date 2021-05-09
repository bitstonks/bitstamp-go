package bitstamp

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsEvent struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type WsClient struct {
	*wsClientConfig
	ws       *websocket.Conn
	done     chan bool
	sendLock sync.Mutex
	Stream   chan *WsEvent
	Errors   chan error
}

func NewWsClient(options ...WsOption) (*WsClient, error) {
	cfg := defaultWsClientConfig()
	for _, opt := range options {
		opt(cfg)
	}

	c := WsClient{
		wsClientConfig: cfg,
		done:           make(chan bool, 1),
		Stream:         make(chan *WsEvent),
		Errors:         make(chan error),
	}

	// set up websocket
	ws, _, err := websocket.DefaultDialer.Dial(c.domain, nil)
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
			c.ws.SetReadDeadline(time.Now().Add(c.timeout))
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

// Determines whether server is requesting reconnect. If such a request is made by the server,
// we should immediately reconnect.
// Note: Bitstamp ensures, that once such a request is received by the client, any new websocket client is connected
// to a healthy server.
func (c *WsClient) IsReconnectRequest(event *WsEvent) bool {
	return event.Event == "bts:request_reconnect"
}

func (c *WsClient) sendEvent(sub WsEvent) {
	c.sendLock.Lock()
	defer c.sendLock.Unlock()

	err := c.ws.WriteJSON(&sub)
	if err != nil {
		fmt.Println(err)
	}
}
