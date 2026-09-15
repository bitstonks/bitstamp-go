package http

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const bitstampHttpApiUrl = "https://www.bitstamp.net/api"

type httpClientConfig struct {
	client             *http.Client
	domain             url.URL
	apiKey             string
	apiSecret          string
	nonceGenerator     func() string
	timestampGenerator func() string
	// have client implicitly round input prices/amounts to correct number of decimal places.
	// used solely for consumers' convenience and will probably be removed at some point.
	autoRounding   bool
	requestTimeout time.Duration
}

func defaultHttpClientConfig() *httpClientConfig {
	domain, err := url.Parse(bitstampHttpApiUrl)
	if err != nil {
		log.Panicf("error parsing domain %s: %v", bitstampHttpApiUrl, err)
	}
	return &httpClientConfig{
		client:             http.DefaultClient,
		domain:             *domain,
		nonceGenerator:     defaultNonce,
		timestampGenerator: timestamp,
		requestTimeout:     10 * time.Second,
	}
}

type HttpOption func(*httpClientConfig)

func UrlDomain(rawDomain string) HttpOption {
	domain, err := url.Parse(rawDomain)
	if err != nil {
		log.Panicf("error parsing domain %s: %v", rawDomain, err)
	}
	return func(config *httpClientConfig) {
		config.domain = *domain
	}
}

func Credentials(apiKey string, apiSecret string) HttpOption {
	return func(config *httpClientConfig) {
		config.apiKey = apiKey
		config.apiSecret = apiSecret
	}
}

func AutoRoundingEnabled() HttpOption {
	return func(config *httpClientConfig) {
		config.autoRounding = true
	}
}

func WithHttpClient(cl *http.Client) HttpOption {
	return func(config *httpClientConfig) {
		config.client = cl
	}
}

func RequestTimeout(t time.Duration) HttpOption {
	return func(config *httpClientConfig) {
		config.requestTimeout = t
	}
}

// 10x slower the than previous `fmt.Sprintf("%d", time.Now().UnixNano())`, should I worry?
func defaultNonce() string {
	return uuid.NewString()
}

// not having an HttpOption for this by design, the timestamp is prescribed by API specification
func timestamp() string {
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano()/1000000)
}
