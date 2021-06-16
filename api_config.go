package bitstamp

import (
	"log"
	"net/url"

	"github.com/google/uuid"
)

const bitstampApiUrl = "https://www.bitstamp.net/api"

type apiClientConfig struct {
	domain         url.URL
	username       string
	apiKey         string
	apiSecret      string
	nonceGenerator func() string
}

func defaultApiClientConfig() *apiClientConfig {
	domain, err := url.Parse(bitstampApiUrl)
	if err != nil {
		log.Panicf("error parsing domain %s: %v", bitstampApiUrl, err)
	}
	return &apiClientConfig{
		domain:         *domain,
		nonceGenerator: defaultNonce,
	}
}

type ApiOption func(*apiClientConfig)

func UrlDomain(rawDomain string) ApiOption {
	domain, err := url.Parse(rawDomain)
	if err != nil {
		log.Panicf("error parsing domain %s: %v", rawDomain, err)
	}
	return func(config *apiClientConfig) {
		config.domain = *domain
	}
}

func Credentials(customerId string, apiKey string, apiSecret string) ApiOption {
	return func(config *apiClientConfig) {
		config.username = customerId
		config.apiKey = apiKey
		config.apiSecret = apiSecret
	}
}

func NonceGenerator(nonceGen func() string) ApiOption {
	return func(config *apiClientConfig) {
		config.nonceGenerator = nonceGen
	}
}

// 10x slower the than previous `fmt.Sprintf("%d", time.Now().UnixNano())`, should i worry?
func defaultNonce() string {
	return uuid.NewString()
}
