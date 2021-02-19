package bitstamp

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlMerge(t *testing.T) {
	bitstampUrl, _ := url.Parse(bitstampApiUrl)
	contrivedUrl, _ := url.Parse("http://127.0.0.1:9876")
	cases := []struct {
		urlBase        url.URL
		path           string
		queryParams    [][2]string
		expectedResult string
	}{
		{*bitstampUrl, "asdf", [][2]string{}, "https://www.bitstamp.net/api/asdf"},
		{*bitstampUrl, "", [][2]string{{"q", "1"}}, "https://www.bitstamp.net/api?q=1"},
		{*bitstampUrl, "v2/ticker", [][2]string{}, "https://www.bitstamp.net/api/v2/ticker"},
		{*bitstampUrl, "v2/ticker/", [][2]string{}, "https://www.bitstamp.net/api/v2/ticker/"},
		{*bitstampUrl, "/v2/ticker/", [][2]string{{"q", "3"}, {"t", "asdf"}}, "https://www.bitstamp.net/api/v2/ticker/?q=3&t=asdf"},
		{*contrivedUrl, "api/v2/ticker/", [][2]string{{"q", "3"}, {"t", "asdf"}}, "http://127.0.0.1:9876/api/v2/ticker/?q=3&t=asdf"},
	}

	for _, c := range cases {
		t.Run("test url merge", func(t *testing.T) {
			actual := urlMerge(c.urlBase, c.path, c.queryParams...)
			assert.Equal(t, c.expectedResult, actual)
		})
	}
}
