package http

import (
	"net/url"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestUrlMerge(t *testing.T) {
	bitstampUrl, _ := url.Parse(bitstampHttpApiUrl)
	contrivedUrl, _ := url.Parse("http://127.0.0.1:9876")
	cases := []struct {
		urlBase        url.URL
		path           string
		queryParams    url.Values
		expectedResult string
	}{
		{*bitstampUrl, "asdf", url.Values{}, "https://www.bitstamp.net/api/asdf"},
		{*bitstampUrl, "", map[string][]string{"q": {"1"}}, "https://www.bitstamp.net/api?q=1"},
		{*bitstampUrl, "v2/ticker", url.Values{}, "https://www.bitstamp.net/api/v2/ticker"},
		{*bitstampUrl, "v2/ticker/", url.Values{}, "https://www.bitstamp.net/api/v2/ticker/"},
		{*bitstampUrl, "/v2/ticker/", map[string][]string{"q": {"3"}, "t": {"asdf"}}, "https://www.bitstamp.net/api/v2/ticker/?q=3&t=asdf"},
		{*contrivedUrl, "api/v2/ticker/", map[string][]string{"q": {"3"}, "t": {"asdf"}}, "http://127.0.0.1:9876/api/v2/ticker/?q=3&t=asdf"},
	}

	for _, c := range cases {
		t.Run("test url merge", func(t *testing.T) {
			actual := urlMerge(c.urlBase, c.path, &c.queryParams)
			assert.Equal(t, c.expectedResult, actual)
		})
	}
}

func TestApiClient_V2Ticker(t *testing.T) {
	c := NewHttpClient()
	resp, err := c.V2Ticker("btcusd")

	assert.NoError(t, err)
	assert.IsType(t, decimal.Decimal{}, resp.Volume)
	assert.True(t, resp.High.GreaterThanOrEqual(resp.Low))
}
