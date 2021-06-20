package bitstamp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

	// apparently, path.Join loses trailing slash in urlPath. we don't want that...
	if strings.HasSuffix(urlPath, "/") {
		baseUrl.Path += "/"
	}

	// add query params
	values := baseUrl.Query()
	for _, param := range queryParams {
		values.Set(param[0], param[1])
	}
	baseUrl.RawQuery = values.Encode()

	return baseUrl.String()
}

type ApiClient struct {
	*apiClientConfig
}

func NewApiClient(options ...ApiOption) *ApiClient {
	config := defaultApiClientConfig()
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

func (c *ApiClient) authenticatedPostRequest(responseObject interface{}, urlPath string, queryParams ...[2]string) (err error) {
	authVersion := "v2"
	method := "POST"
	xAuth := "BITSTAMP " + c.apiKey
	apiSecret := []byte(c.apiSecret)
	contentType := "application/x-www-form-urlencoded"
	timestamp_ := c.timestampGenerator()
	nonce := c.nonceGenerator()
	url_ := urlMerge(c.domain, urlPath)

	var payloadString string
	if queryParams != nil {
		urlParams := url.Values{}
		for _, p := range queryParams {
			urlParams.Set(p[0], p[1]) // TODO: or is it .Add() here? any array arguments in the documentation?
		}
		payloadString = urlParams.Encode()
	}

	// message construction
	msg := xAuth + method + strings.TrimPrefix(url_, "https://")
	if queryParams == nil {
		msg = msg + nonce + timestamp_ + authVersion // TODO: apparently, contentType must be omitted here?
	} else {
		msg = msg + contentType + nonce + timestamp_ + authVersion + payloadString
	}
	sig := hmac.New(sha256.New, apiSecret)
	sig.Write([]byte(msg))
	signature := hex.EncodeToString(sig.Sum(nil))

	// do the request
	client := &http.Client{}
	var req *http.Request
	if queryParams == nil {
		req, err = http.NewRequest(method, url_, nil)
	} else {
		req, err = http.NewRequest(method, url_, bytes.NewBuffer([]byte(payloadString)))
	}
	if err != nil {
		return
	}
	req.Header.Add("X-Auth", xAuth)
	req.Header.Add("X-Auth-Signature", signature)
	req.Header.Add("X-Auth-Nonce", nonce)
	req.Header.Add("X-Auth-Timestamp", timestamp_)
	req.Header.Add("X-Auth-Version", authVersion)
	if queryParams != nil {
		req.Header.Add("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// handle response
	if resp.StatusCode != 200 {
		var errorMsg map[string]interface{}
		err = json.Unmarshal(respBody, &errorMsg)
		if err != nil {
			return
		}

		reasonVal, reasonPresent := errorMsg["reason"]
		codeVal, codePresent := errorMsg["code"]
		if reasonPresent && codePresent {
			err = fmt.Errorf("%s %s (%d)", codeVal, reasonVal, resp.StatusCode)
			return
		} else {
			err = fmt.Errorf("%s (%d)", string(respBody), resp.StatusCode)
			return
		}
	} else {
		// verify server signature
		checkMsg := nonce + timestamp_ + resp.Header.Get("Content-Type") + string(respBody)
		sig := hmac.New(sha256.New, apiSecret)
		sig.Write([]byte(checkMsg))
		serverSig := hex.EncodeToString(sig.Sum(nil))
		if serverSig != resp.Header.Get("X-Server-Auth-Signature") {
			err = fmt.Errorf("server signature mismatch: us (%s) them (%s)", serverSig, resp.Header.Get("X-Server-Auth-Signature"))
			return
		}

		err = json.Unmarshal(respBody, responseObject)
		if err != nil {
			return
		}
	}

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
	err = c.getRequest("/ticker/", &response)
	return
}

// GET https://www.bitstamp.net/api/v2/ticker/{currency_pair}/
func (c *ApiClient) V2Ticker(currencyPair string) (response TickerResponse, err error) {
	urlPath := fmt.Sprintf("/v2/ticker/%s/", currencyPair)
	err = c.getRequest(urlPath, &response)
	return
}

// GET https://www.bitstamp.net/api/ticker_hour/
func (c *ApiClient) V1HourlyTicker() (response TickerResponse, err error) {
	err = c.getRequest("/ticker_hour/", &response)
	return
}

// GET https://www.bitstamp.net/api/v2/ticker_hour/{currency_pair}/
func (c *ApiClient) V2HourlyTicker(currencyPair string) (response TickerResponse, err error) {
	urlPath := fmt.Sprintf("/v2/ticker_hour/%s/", currencyPair)
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
	err = c.getRequest("/order_book/", &response, [2]string{"group", strconv.Itoa(group)})
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
	urlPath := fmt.Sprintf("/v2/order_book/%s/", currencyPair)
	err = c.getRequest(urlPath, &response, [2]string{"group", strconv.Itoa(group)})
	return
}

//
// Private Functions
//

// Account balance
type V2BalanceResponse struct {
	// currencies
	BchAvailable      decimal.Decimal `json:"bch_available"`
	BchBalance        decimal.Decimal `json:"bch_balance"`
	BchReserved       decimal.Decimal `json:"bch_reserved"`
	BchWithdrawalFee  decimal.Decimal `json:"bch_withdrawal_fee"`
	BtcAvailable      decimal.Decimal `json:"btc_available"`
	BtcBalance        decimal.Decimal `json:"btc_balance"`
	BtcReserved       decimal.Decimal `json:"btc_reserved"`
	BtcWithdrawalFee  decimal.Decimal `json:"btc_withdrawal_fee"`
	EthAvailable      decimal.Decimal `json:"eth_available"`
	EthBalance        decimal.Decimal `json:"eth_balance"`
	EthReserved       decimal.Decimal `json:"eth_reserved"`
	EthWithdrawalFee  decimal.Decimal `json:"eth_withdrawal_fee"`
	EurAvailable      decimal.Decimal `json:"eur_available"`
	EurBalance        decimal.Decimal `json:"eur_balance"`
	EurReserved       decimal.Decimal `json:"eur_reserved"`
	EurWithdrawalFee  decimal.Decimal `json:"eur_withdrawal_fee"`
	GbpAvailable      decimal.Decimal `json:"gbp_available"`
	GbpBalance        decimal.Decimal `json:"gbp_balance"`
	GbpReserved       decimal.Decimal `json:"gbp_reserved"`
	GbpWithdrawalFee  decimal.Decimal `json:"gbp_withdrawal_fee"`
	LinkAvailable     decimal.Decimal `json:"link_available"`
	LinkBalance       decimal.Decimal `json:"link_balance"`
	LinkReserved      decimal.Decimal `json:"link_reserved"`
	LinkWithdrawalFee decimal.Decimal `json:"link_withdrawal_fee"`
	LtcAvailable      decimal.Decimal `json:"ltc_available"`
	LtcBalance        decimal.Decimal `json:"ltc_balance"`
	LtcReserved       decimal.Decimal `json:"ltc_reserved"`
	LtcWithdrawalFee  decimal.Decimal `json:"ltc_withdrawal_fee"`
	OmgAvailable      decimal.Decimal `json:"omg_available"`
	OmgBalance        decimal.Decimal `json:"omg_balance"`
	OmgReserved       decimal.Decimal `json:"omg_reserved"`
	OmgWithdrawalFee  decimal.Decimal `json:"omg_withdrawal_fee"`
	PaxAvailable      decimal.Decimal `json:"pax_available"`
	PaxBalance        decimal.Decimal `json:"pax_balance"`
	PaxReserved       decimal.Decimal `json:"pax_reserved"`
	PaxWithdrawalFee  decimal.Decimal `json:"pax_withdrawal_fee"`
	UsdAvailable      decimal.Decimal `json:"usd_available"`
	UsdBalance        decimal.Decimal `json:"usd_balance"`
	UsdReserved       decimal.Decimal `json:"usd_reserved"`
	UsdWithdrawalFee  decimal.Decimal `json:"usd_withdrawal_fee"`
	UsdcAvailable     decimal.Decimal `json:"usdc_available"`
	UsdcBalance       decimal.Decimal `json:"usdc_balance"`
	UsdcReserved      decimal.Decimal `json:"usdc_reserved"`
	UsdcWithdrawalFee decimal.Decimal `json:"usdc_withdrawal_fee"`
	XlmAvailable      decimal.Decimal `json:"xlm_available"`
	XlmBalance        decimal.Decimal `json:"xlm_balance"`
	XlmReserved       decimal.Decimal `json:"xlm_reserved"`
	XlmWithdrawalFee  decimal.Decimal `json:"xlm_withdrawal_fee"`
	XrpAvailable      decimal.Decimal `json:"xrp_available"`
	XrpBalance        decimal.Decimal `json:"xrp_balance"`
	XrpReserved       decimal.Decimal `json:"xrp_reserved"`
	XrpWithdrawalFee  decimal.Decimal `json:"xrp_withdrawal_fee"`

	// pairs
	BchbtcFee  decimal.Decimal `json:"bchbtc_fee"`
	BcheurFee  decimal.Decimal `json:"bcheur_fee"`
	BchgbpFee  decimal.Decimal `json:"bchgbp_fee"`
	BchusdFee  decimal.Decimal `json:"bchusd_fee"`
	BtceurFee  decimal.Decimal `json:"btceur_fee"`
	BtcgbpFee  decimal.Decimal `json:"btcgbp_fee"`
	BtcpaxFee  decimal.Decimal `json:"btcpax_fee"`
	BtcusdcFee decimal.Decimal `json:"btcusdc_fee"`
	BtcusdFee  decimal.Decimal `json:"btcusd_fee"`
	EthbtcFee  decimal.Decimal `json:"ethbtc_fee"`
	EtheurFee  decimal.Decimal `json:"etheur_fee"`
	EthgbpFee  decimal.Decimal `json:"ethgbp_fee"`
	EthpaxFee  decimal.Decimal `json:"ethpax_fee"`
	EthusdcFee decimal.Decimal `json:"ethusdc_fee"`
	EthusdFee  decimal.Decimal `json:"ethusd_fee"`
	EurusdFee  decimal.Decimal `json:"eurusd_fee"`
	GbpeurFee  decimal.Decimal `json:"gbpeur_fee"`
	GbpusdFee  decimal.Decimal `json:"gbpusd_fee"`
	LinkbtcFee decimal.Decimal `json:"linkbtc_fee"`
	LinkethFee decimal.Decimal `json:"linketh_fee"`
	LinkeurFee decimal.Decimal `json:"linkeur_fee"`
	LinkgbpFee decimal.Decimal `json:"linkgbp_fee"`
	LinkusdFee decimal.Decimal `json:"linkusd_fee"`
	LtcbtcFee  decimal.Decimal `json:"ltcbtc_fee"`
	LtceurFee  decimal.Decimal `json:"ltceur_fee"`
	LtcgbpFee  decimal.Decimal `json:"ltcgbp_fee"`
	LtcusdFee  decimal.Decimal `json:"ltcusd_fee"`
	OmgbtcFee  decimal.Decimal `json:"omgbtc_fee"`
	OmgeurFee  decimal.Decimal `json:"omgeur_fee"`
	OmggbpFee  decimal.Decimal `json:"omggbp_fee"`
	OmgusdFee  decimal.Decimal `json:"omgusd_fee"`
	PaxeurFee  decimal.Decimal `json:"paxeur_fee"`
	PaxgbpFee  decimal.Decimal `json:"paxgbp_fee"`
	PaxusdFee  decimal.Decimal `json:"paxusd_fee"`
	UsdceurFee decimal.Decimal `json:"usdceur_fee"`
	UsdcusdFee decimal.Decimal `json:"usdcusd_fee"`
	XlmbtcFee  decimal.Decimal `json:"xlmbtc_fee"`
	XlmeurFee  decimal.Decimal `json:"xlmeur_fee"`
	XlmgbpFee  decimal.Decimal `json:"xlmgbp_fee"`
	XlmusdFee  decimal.Decimal `json:"xlmusd_fee"`
	XrpbtcFee  decimal.Decimal `json:"xrpbtc_fee"`
	XrpeurFee  decimal.Decimal `json:"xrpeur_fee"`
	XrpgbpFee  decimal.Decimal `json:"xrpgbp_fee"`
	XrppaxFee  decimal.Decimal `json:"xrppax_fee"`
	XrpusdFee  decimal.Decimal `json:"xrpusd_fee"`

	// fee
	Fee decimal.Decimal `json:"fee"`
}

// POST https://www.bitstamp.net/api/v2/balance/
// POST https://www.bitstamp.net/api/v2/balance/{currency_pair}/
func (c *ApiClient) V2Balance(currencyPairOrAll string) (response V2BalanceResponse, err error) {
	// TODO: validate currency pair
	if currencyPairOrAll == "all" {
		err = c.authenticatedPostRequest(&response, "/v2/balance/")
	} else {
		err = c.authenticatedPostRequest(&response, fmt.Sprintf("/v2/balance/%s/", currencyPairOrAll))
	}

	return
}

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
		urlPath = "/v2/open_orders/all/"
	} else {
		urlPath = fmt.Sprintf("/v2/open_orders/%s/", currencyPairOrAll)
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
	err = c.authenticatedPostRequest(&response, "/v2/cancel_order/", [2]string{"id", orderId})
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
	urlPath := fmt.Sprintf("/v2/%s/%s/", side, currencyPair)
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

type V2MarketOrderResponse struct {
	Id       string          `json:"id"`
	Datetime string          `json:"datetime"`
	Type     string          `json:"type"`
	Price    decimal.Decimal `json:"price"`
	Amount   decimal.Decimal `json:"amount"`
	Error    string          `json:"error"`
	Status   string          `json:"status"`
	Reason   interface{}     `json:"reason"`
}

func (c *ApiClient) v2MarketOrder(side, currencyPair string, amount decimal.Decimal, clOrdId string) (response V2MarketOrderResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s/market/%s/", side, currencyPair)
	url_ := urlMerge(c.domain, urlPath)

	data := c.credentials()
	data.Set("amount", amount.String())
	if clOrdId != "" {
		data.Set("client_order_id", clOrdId)
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
		err = fmt.Errorf("error placing market %s (for %s): %v", side, amount, response.Reason)
		return
	}

	return
}

func (c *ApiClient) V2BuyMarketOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2MarketOrderResponse, err error) {
	return c.v2MarketOrder("buy", currencyPair, amount, clOrdId)
}

func (c *ApiClient) V2SellMarketOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2MarketOrderResponse, err error) {
	return c.v2MarketOrder("sell", currencyPair, amount, clOrdId)
}

type V2InstantOrderResponse struct {
	Id       string          `json:"id"`
	Datetime string          `json:"datetime"`
	Type     string          `json:"type"`
	Price    decimal.Decimal `json:"price"`
	Amount   decimal.Decimal `json:"amount"`
	Error    string          `json:"error"`
	Status   string          `json:"status"`
	Reason   interface{}     `json:"reason"`
}

func (c *ApiClient) v2InstantOrder(side, currencyPair string, amount decimal.Decimal, clOrdId string) (response V2InstantOrderResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s/instant/%s/", side, currencyPair)
	url_ := urlMerge(c.domain, urlPath)

	data := c.credentials()
	data.Set("amount", amount.String())
	if clOrdId != "" {
		data.Set("client_order_id", clOrdId)
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
		err = fmt.Errorf("error placing instant %s (for %s): %v", side, amount, response.Reason)
		return
	}

	return
}

func (c *ApiClient) V2BuyInstantOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2InstantOrderResponse, err error) {
	return c.v2InstantOrder("buy", currencyPair, amount, clOrdId)
}

func (c *ApiClient) V2SellInstantOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2InstantOrderResponse, err error) {
	return c.v2InstantOrder("sell", currencyPair, amount, clOrdId)
}
