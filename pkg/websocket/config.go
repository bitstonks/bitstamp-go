package websocket

import (
	"time"
)

const bitstampWsUrl = "wss://ws.bitstamp.net"
const wsTimeout = 60 * time.Second

type wsClientConfig struct {
	domain  string
	timeout time.Duration
}

func defaultWsClientConfig() *wsClientConfig {
	return &wsClientConfig{
		domain:  bitstampWsUrl,
		timeout: wsTimeout,
	}
}

type WsOption func(*wsClientConfig)

func WsUrl(domain string) WsOption {
	return func(config *wsClientConfig) {
		config.domain = domain
	}
}

func Timeout(timeout time.Duration) WsOption {
	return func(config *wsClientConfig) {
		config.timeout = timeout
	}
}
