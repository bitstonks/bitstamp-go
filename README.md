[![Go Report Card](https://goreportcard.com/badge/github.com/bitstonks/bitstamp-go)](https://goreportcard.com/report/github.com/bitstonks/bitstamp-go)

# bitstamp-go

Client implementations for Bitstamp's HTTP (REST) and Websocket APIs in Go. Websocket client uses
[API v2](https://www.bitstamp.net/websocket/v2/) instead of Pusher (which is not available anymore).

## Requirements

* Go 1.18+
* Dependencies in `go.mod`

## Usage examples

Require the repo in your `go.mod` file. See examples folder for more. For instance, you can try running:

```bash
$ go run examples/http/main.go
```

## TODO

* Configure GitHub Actions.
* Update pairs config.
* Finish implementing all the endpoints.
* Remove/deprecate REST v1.
* Godoc / documentation.
* More tests.
* E2E tests against Bitstamp's API.
* Docker builds.
* Authors file / information?
* Shadow order book creation. (Is this an example? A separate project even?)
