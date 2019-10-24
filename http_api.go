package bitstamp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func urlMerge(baseUrl url.URL, urlPath string, queryParams ...[2]string) string {
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)

	// add query params
	values := baseUrl.Query()
	for _, param := range queryParams {
		values.Set(param[0], param[1])
	}
	baseUrl.RawQuery = values.Encode()

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

// TODO: change the order of method arguments here...
func (c *ApiClient) getRequest(urlPath string, responseObject interface{}, queryParams ...[2]string) (err error) {
	url_ := urlMerge(c.domain, urlPath, queryParams...)

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

type OrderBookEntry struct {
	Price  decimal.Decimal
	Amount decimal.Decimal
	Id     uint64
}

// helper method that lets us unmarshal order book entries directly from API JSON responses
func (obe *OrderBookEntry) UnmarshalJSON(bytes []byte) error {
	if obe == nil {
		obe = new(OrderBookEntry)
	}

	var parts []string
	err := json.Unmarshal(bytes, &parts)
	if err != nil {
		return err
	}

	if len(parts) != 2 && len(parts) != 3 {
		return fmt.Errorf("wrong number of arguments for PriceAmountId: %v", parts)
	}

	price, err := decimal.NewFromString(parts[0])
	if err != nil {
		return err
	}
	obe.Price = price

	amount, err := decimal.NewFromString(parts[1])
	if err != nil {
		return err
	}
	obe.Amount = amount

	if len(parts) == 3 {
		orderId, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}

		obe.Id = uint64(orderId)
	}

	return nil
}

type V1OrderBookResponse struct {
	Timestamp string           `json:"timestamp"` // UNIX epoch in UTC in seconds
	Bids      []OrderBookEntry `json:"bids"`
	Asks      []OrderBookEntry `json:"asks"`
}

// GET https://www.bitstamp.net/api/order_book?group=1
func (c *ApiClient) V1OrderBook(group int) (response V1OrderBookResponse, err error) {
	err = c.getRequest("/api/order_book/", &response, [2]string{"group", strconv.Itoa(group)})
	return
}

type V2OrderBookResponse struct {
	V1OrderBookResponse
	Microtimestamp string `json:"microtimestamp"`
}

// GET https://www.bitstamp.net/api/v2/order_book/{currency_pair}?group=1
// Possible values are for group parameter
// - 0 (orders are not grouped at same price)
// - 1 (orders are grouped at same price - default)
// - 2 (orders with their order ids are not grouped at same price)
// GET https://www.bitstamp.net/api/order_book?group=1
func (c *ApiClient) V2OrderBook(currencyPair string, group int) (response V2OrderBookResponse, err error) {
	urlPath := fmt.Sprintf("/api/v2/order_book/%s/", currencyPair)
	err = c.getRequest(urlPath, &response, [2]string{"group", strconv.Itoa(group)})
	return
}

//
// Private Functions
//

// Account balance
// User transactions
// Open orders
type V2OpenOrdersResponse struct {
	Id           string          `json:"id"`
	Datetime     string          `json:"datetime"`
	Type         string          `json:"type"`
	Price        decimal.Decimal `json:"price"`
	Amount       decimal.Decimal `json:"amount"`
	CurrencyPair string          `json:"currency_pair"`
	Status       string          `json:"status"`
	Reason       interface{}     `json:"reason"`
}

// POST https://www.bitstamp.net/api/v2/open_orders/all/
// POST https://www.bitstamp.net/api/v2/open_orders/{currency_pair}
func (c *ApiClient) V2OpenOrders(currencyPairOrAll string) (response []V2OpenOrdersResponse, err error) {
	var urlPath string
	if currencyPairOrAll == "all" {
		urlPath = "/api/v2/open_orders/all/"
	} else {
		urlPath = fmt.Sprintf("/api/v2/open_orders/%s/", currencyPairOrAll)
	}
	url_ := urlMerge(c.domain, urlPath)

	resp, err := http.PostForm(url_, c.credentials())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	return
}

// Order status
// Cancel order
type V2CancelOrderResponse struct {
	Id     uint64          `json:"id"`
	Amount decimal.Decimal `json:"amount"`
	Price  decimal.Decimal `json:"price"`
	Type   uint8           `json:"type"`
	Error  string          `json:"error"`
}

func (c *ApiClient) V2CancelOrder(orderId string) (response V2CancelOrderResponse, err error) {
	url_ := urlMerge(c.domain, "/api/v2/cancel_order/")

	data := c.credentials()
	data.Set("id", orderId)

	resp, err := http.PostForm(url_, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.Error != "" {
		err = errors.New(response.Error)
		return
	}

	return
}

// Cancel all orders
// Buy limit order
// Sell limit order

//{"status": "error", "reason": {"__all__": ["Price is more than 20% below market price."]}}
//{"status": "error", "reason": {"__all__": ["You need 158338.86 USD to open that order. You have only 99991.52 USD available. Check your account balance for details."]}}
type V2LimitOrderResponse struct {
	Id       string          `json:"id"`
	Datetime string          `json:"datetime"`
	Type     string          `json:"type"`
	Price    decimal.Decimal `json:"price"`
	Amount   decimal.Decimal `json:"amount"`
	Status   string          `json:"status"`
	Reason   interface{}     `json:"reason"`
}

func (c *ApiClient) v2LimitOrder(side, currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
	urlPath := fmt.Sprintf("/api/v2/%s/%s/", side, currencyPair)
	url_ := urlMerge(c.domain, urlPath)

	data := c.credentials()
	data.Set("price", price.String())
	data.Set("amount", amount.String())
	if clOrdId != "" {
		data.Set("client_order_id", clOrdId)
	}
	if dailyOrder {
		data.Set("daily_order", "True")
	}
	if iocOrder {
		data.Set("ioc_order", "True")
	}

	resp, err := http.PostForm(url_, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error placing limit %s (%s @ %s): %v", side, amount, price, response.Reason)
		return
	}

	return
}

func (c *ApiClient) V2BuyLimitOrder(currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
	return c.v2LimitOrder("buy", currencyPair, price, amount, limitPrice, dailyOrder, iocOrder, clOrdId)
}

func (c *ApiClient) V2SellLimitOrder(currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
	return c.v2LimitOrder("sell", currencyPair, price, amount, limitPrice, dailyOrder, iocOrder, clOrdId)
}
