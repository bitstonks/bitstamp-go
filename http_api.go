package bitstamp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/shopspring/decimal"
)

func urlMerge(baseUrl url.URL, urlPath string) string {
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)
	return baseUrl.String()
}

type ApiClient struct {
	*clientConfig
}

func NewApiClient(options ...Option) *ApiClient {
	config := defaultClientConfig()
	for _, option := range options {
		option(config)
	}
	return &ApiClient{config}
}

func (c *ApiClient) credentials() url.Values {
	nonce := c.nonceGenerator()
	message := nonce + c.username + c.apiKey

	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(message))
	signature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	data := make(url.Values)
	data.Set("key", c.apiKey)
	data.Set("signature", signature)
	data.Set("nonce", nonce)
	return data
}

func (c *ApiClient) getRequest(urlPath string, responseObject interface{}) (err error) {
	url_ := urlMerge(c.domain, urlPath)

	resp, err := http.Get(url_)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, responseObject)
	return
}

//
// Public methods
//

type TickerResponse struct {
	High      decimal.Decimal `json:"high"`
	Last      decimal.Decimal `json:"last"`
	Timestamp decimal.Decimal `json:"timestamp"`
	Bid       decimal.Decimal `json:"bid"`
	Vwap      decimal.Decimal `json:"vwap"`
	Volume    decimal.Decimal `json:"volume"`
	Low       decimal.Decimal `json:"low"`
	Ask       decimal.Decimal `json:"ask"`
	Open      decimal.Decimal `json:"open"`
}

// GET https://www.bitstamp.net/api/ticker/
func (c *ApiClient) V1Ticker() (response TickerResponse, err error) {
	err = c.getRequest("/api/ticker/", &response)
	return
}

// GET https://www.bitstamp.net/api/v2/ticker/{currency_pair}/
func (c *ApiClient) V2Ticker(currencyPair string) (response TickerResponse, err error) {
	urlPath := fmt.Sprintf("/api/v2/ticker/%s/", currencyPair)
	err = c.getRequest(urlPath, &response)
	return
}

// GET https://www.bitstamp.net/api/ticker_hour/
func (c *ApiClient) V1HourlyTicker() (response TickerResponse, err error) {
	err = c.getRequest("/api/ticker_hour/", &response)
	return
}

// GET https://www.bitstamp.net/api/v2/ticker_hour/{currency_pair}/
func (c *ApiClient) V2HourlyTicker(currencyPair string) (response TickerResponse, err error) {
	urlPath := fmt.Sprintf("/api/v2/ticker_hour/%s/", currencyPair)
	err = c.getRequest(urlPath, &response)
	return
}

type V1OrderBookResponse struct {
	Timestamp string               `json:"timestamp"` // UNIX epoch in UTC in seconds
	Bids      [][2]decimal.Decimal `json:"bids"`
	Asks      [][2]decimal.Decimal `json:"asks"`
}

// GET https://www.bitstamp.net/api/order_book/
func (c *ApiClient) V1OrderBook() (response V1OrderBookResponse, err error) {
	err = c.getRequest("/api/order_book/", &response)
	return
}

type ApiV2BalancesResponse struct {
	BchAvailable decimal.Decimal `json:"bch_available"`
	BchBalance   decimal.Decimal `json:"bch_balance"`
	BchReserved  decimal.Decimal `json:"bch_reserved"`
	BchBtcFee    decimal.Decimal `json:"bchbtc_fee"`
	BchEurFee    decimal.Decimal `json:"bcheur_fee"`
	BchUsdFee    decimal.Decimal `json:"bchusd_fee"`

	BtcAvailable decimal.Decimal `json:"btc_available"`
	BtcBalance   decimal.Decimal `json:"btc_balance"`
	BtcReserved  decimal.Decimal `json:"btc_reserved"`
	BtcEurFee    decimal.Decimal `json:"btceur_fee"`
	BtcUsdFee    decimal.Decimal `json:"btcusd_fee"`

	EthAvailable decimal.Decimal `json:"eth_available"`
	EthBalance   decimal.Decimal `json:"eth_balance"`
	EthReserved  decimal.Decimal `json:"eth_reserved"`
	EthBtcFee    decimal.Decimal `json:"ethbtc_fee"`
	EthEurFee    decimal.Decimal `json:"etheur_fee"`
	EthUsdFee    decimal.Decimal `json:"ethusd_fee"`

	EurAvailable decimal.Decimal `json:"eur_available"`
	EurBalance   decimal.Decimal `json:"eur_balance"`
	EurReserved  decimal.Decimal `json:"eur_reserved"`
	EurUsdFee    decimal.Decimal `json:"eurusd_fee"`

	LtcAvailable decimal.Decimal `json:"ltc_available"`
	LtcBalance   decimal.Decimal `json:"ltc_balance"`
	LtcReserved  decimal.Decimal `json:"ltc_reserved"`
	LtcBtcFee    decimal.Decimal `json:"ltcbtc_fee"`
	LtcEurFee    decimal.Decimal `json:"ltceur_fee"`
	LtcUsdFee    decimal.Decimal `json:"ltcusd_fee"`

	UsdAvailable decimal.Decimal `json:"usd_available"`
	UsdBalance   decimal.Decimal `json:"usd_balance"`
	UsdReserved  decimal.Decimal `json:"usd_reserved"`

	XrpAvailable decimal.Decimal `json:"xrp_available"`
	XrpBalance   decimal.Decimal `json:"xrp_balance"`
	XrpReserved  decimal.Decimal `json:"xrp_reserved"`
	XrpBtcFee    decimal.Decimal `json:"xrpbtc_fee"`
	XrpEurFee    decimal.Decimal `json:"xrpeur_fee"`
	XrpUsdFee    decimal.Decimal `json:"xrpusd_fee"`
}

func (c *ApiClient) ApiV2Balances() (response ApiV2BalancesResponse, err error) {
	endpoint := fmt.Sprintf("%s/api/v2/balance/", c.domain)

	resp, err := http.PostForm(endpoint, c.credentials())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &response)
	return
}
