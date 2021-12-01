[![Go Report Card](https://goreportcard.com/badge/github.com/samotarnik/bitstamp-go)](https://goreportcard.com/report/github.com/samotarnik/bitstamp-go)

# bitstamp-go

Client implementations for Bitstamp's REST and Websocket APIs in Go. Copied
heavily from [github.com/ajph/bitstamp-go](https://github.com/ajph/bitstamp-go).
Websocket client uses [API v2](https://www.bitstamp.net/websocket/v2/) instead
of Pusher (which is not available anymore).

## Requirements

* Go 1.15+
* Dependencies in `go.mod`

## Usage examples

Require the repo in your `go.mod` file. See examples folder for more. For
instance, you can try running:

```bash
$ go run examples/http_api/http_api.go
```

## TODO

* Upgrade Go.
* Configure Github Actions.
* Update pairs.
* Finish implementing all the endpoints.
* Remove/deprecate REST v1.
* Godoc / documentation.
* More tests.
* E2E tests against Bitstamp's API.
* Docker builds.
* Authors file / information?
* Shadow order book creation. (Is this an example? A separate project even?)
