package bitstamp

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

const bitstampApiUrl = "https://www.bitstamp.net"

type clientConfig struct {
	domain         url.URL
	username       string
	apiKey         string
	apiSecret      string
	nonceGenerator func() string
}

func defaultClientConfig() *clientConfig {
	domain, err := url.Parse(bitstampApiUrl)
	if err != nil {
		log.Panicf("error parsing domain %s: %v", bitstampApiUrl, err)
	}
	return &clientConfig{
		domain:         *domain,
		nonceGenerator: defaultNonce,
	}
}

type Option func(*clientConfig)

func UrlDomain(rawDomain string) Option {
	domain, err := url.Parse(rawDomain)
	if err != nil {
		log.Panicf("error parsing domain %s: %v", rawDomain, err)
	}
	return func(config *clientConfig) {
		config.domain = *domain
	}
}

func Credentials(customerId string, apiKey string, apiSecret string) Option {
	return func(config *clientConfig) {
		config.username = customerId
		config.apiKey = apiKey
		config.apiSecret = apiSecret
	}
}

func NonceGenerator(nonceGen func() string) Option {
	return func(config *clientConfig) {
		config.nonceGenerator = nonceGen
	}
}

func defaultNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
