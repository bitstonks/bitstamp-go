[![Go Report Card](https://goreportcard.com/badge/github.com/samotarnik/bitstamp-go)](https://goreportcard.com/report/github.com/samotarnik/bitstamp-go)
[![CircleCI](https://circleci.com/gh/samotarnik/bitstamp-go/tree/dump.svg?style=svg)](https://circleci.com/gh/samotarnik/bitstamp-go/tree/dump)

# bitstamp-go

Client implementations for Bitstamp's REST and Websocket APIs in Go. Copied heavily from [github.com/ajph/bitstamp-go](https://github.com/ajph/bitstamp-go). Websocket client uses [API v2](https://www.bitstamp.net/websocket/v2/) instead of Pusher.

## Requirements

* Go 1.12+
* Dependencies in `go.mod`

## Usage examples

Require the repo in your `go.mod` file. See examples folder for more. For instance, you can try running 
```bash
$ go run examples/http_api/http_api.go
```

## TODO

* godoc
* tests
* Docker builds