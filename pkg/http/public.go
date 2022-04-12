package http

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
	"strings"
)

// Contains "public" endpoints whereby we are following the naming here: https://www.bitstamp.net/api/

//
// Tickers
//

type TickerResponse struct {
	Ask       decimal.Decimal `json:"ask"`
	Bid       decimal.Decimal `json:"bid"`
	High      decimal.Decimal `json:"high"`
	Last      decimal.Decimal `json:"last"`
	Low       decimal.Decimal `json:"low"`
	Open      decimal.Decimal `json:"open"`
	Timestamp int64           `json:"timestamp"`
	Volume    decimal.Decimal `json:"volume"`
	Vwap      decimal.Decimal `json:"vwap"`
}

// custom deserialization instructions necessary
func (t *TickerResponse) UnmarshalJSON(data []byte) error {
	if t == nil {
		t = new(TickerResponse)
	}

	raw := make(map[string]string)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for _, k := range []string{"Ask", "Bid", "High", "Last", "Low", "Open", "Volume", "Vwap"} {
		kl := strings.ToLower(k)
		if strVal, exists := raw[kl]; exists {
			val, err2 := decimal.NewFromString(strVal)
			if err2 != nil {
				return fmt.Errorf("error extracting %s: %v", kl, err2)
			}

			tv := reflect.ValueOf(t).Elem().FieldByName(k)
			if tv.IsValid() {
				tv.Set(reflect.ValueOf(val))
			}
		}
	}

	if val, exists := raw["timestamp"]; exists {
		tmstmp, err2 := strconv.ParseInt(val, 10, 64)
		if err2 != nil {
			return fmt.Errorf("error parsing timestamp: %v", err2)
		}
		t.Timestamp = tmstmp
	} else {
		return fmt.Errorf("missing expected `timestamp` field in %v", raw)
	}

	return nil
}

// GET https://www.bitstamp.net/api/ticker/
func (c *HttpClient) V1Ticker() (response TickerResponse, err error) {
	err = c.getRequest(&response, "/ticker/")
	return
}

// GET https://www.bitstamp.net/api/ticker_hour/
func (c *HttpClient) V1HourlyTicker() (response TickerResponse, err error) {
	err = c.getRequest(&response, "/ticker_hour/")
	return
}

// GET https://www.bitstamp.net/api/v2/ticker/{currency_pair}/
func (c *HttpClient) V2Ticker(currencyPair string) (response TickerResponse, err error) {
	if err = validateCurrencyPair(currencyPair); err != nil {
		return
	}
	urlPath := fmt.Sprintf("/v2/ticker/%s/", currencyPair)
	err = c.getRequest(&response, urlPath)
	return
}

// GET https://www.bitstamp.net/api/v2/ticker_hour/{currency_pair}/
func (c *HttpClient) V2HourlyTicker(currencyPair string) (response TickerResponse, err error) {
	if err = validateCurrencyPair(currencyPair); err != nil {
		return
	}
	urlPath := fmt.Sprintf("/v2/ticker_hour/%s/", currencyPair)
	err = c.getRequest(&response, urlPath)
	return
}

//
// Order books
//

type OrderBookEntry struct {
	Price  decimal.Decimal
	Amount decimal.Decimal
	Id     int64
}

// custom deserialization instructions necessary
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

		obe.Id = orderId
	}

	return nil
}

type V1OrderBookResponse struct {
	Timestamp string           `json:"timestamp"` // UNIX epoch in UTC in seconds
	Bids      []OrderBookEntry `json:"bids"`
	Asks      []OrderBookEntry `json:"asks"`
}

// GET https://www.bitstamp.net/api/order_book?group=1
func (c *HttpClient) V1OrderBook(group int) (response V1OrderBookResponse, err error) {
	err = c.getRequest(&response, "/order_book/", [2]string{"group", strconv.Itoa(group)})
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
func (c *HttpClient) V2OrderBook(currencyPair string, group int) (response V2OrderBookResponse, err error) {
	if err = validateCurrencyPair(currencyPair); err != nil {
		return
	}
	switch group {
	case 0, 1, 2:
		urlPath := fmt.Sprintf("/v2/order_book/%s/", currencyPair)
		err = c.getRequest(&response, urlPath, [2]string{"group", strconv.Itoa(group)})
	default:
		err = fmt.Errorf("invalid group parameter value: %d", group)
	}
	return
}

//
// Transactions
//

type V2TransactionsResponse struct {
	Amount decimal.Decimal `json:"amount"` // Amount in base (?)
	Date   int64           `json:"date"`   // Unix timestamp date and time.
	Price  decimal.Decimal `json:"price"`
	Tid    int64           `json:"tid"`  // Transaction ID.
	Type   int8            `json:"type"` // 0 (buy) or 1 (sell).
}

// custom deserialization instructions necessary
func (tx *V2TransactionsResponse) UnmarshalJSON(bytes []byte) error {
	if tx == nil {
		tx = new(V2TransactionsResponse)
	}

	var raw map[string]string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	if val, exists := raw["amount"]; exists {
		amount, err2 := decimal.NewFromString(val)
		if err2 != nil {
			return fmt.Errorf("error extracting amount: %v", err2)
		}
		tx.Amount = amount
	} else {
		return fmt.Errorf("missing expected `amount` field in %v", raw)
	}

	if val, exists := raw["date"]; exists {
		date, err2 := strconv.ParseInt(val, 10, 64)
		if err2 != nil {
			return fmt.Errorf("error extracting date: %v", err2)
		}
		tx.Date = date
	} else {
		return fmt.Errorf("missing expected `date` field in %v", raw)
	}

	if val, exists := raw["price"]; exists {
		price, err2 := decimal.NewFromString(val)
		if err2 != nil {
			return fmt.Errorf("error extracting price: %v", err2)
		}
		tx.Price = price
	} else {
		return fmt.Errorf("missing expected `price` field in %v", raw)
	}

	if val, exists := raw["tid"]; exists {
		txId, err2 := strconv.ParseInt(val, 10, 64)
		if err2 != nil {
			return fmt.Errorf("error extracting transaction id: %v", err2)
		}
		tx.Tid = txId
	} else {
		return fmt.Errorf("missing expected `tid` field in %v", raw)
	}

	if val, exists := raw["type"]; exists {
		type_, err2 := strconv.ParseInt(val, 10, 8)
		if err2 != nil {
			return fmt.Errorf("error extracting type: %v", err2)
		}
		tx.Type = int8(type_)
	} else {
		return fmt.Errorf("missing expected `tid` field in %v", raw)
	}

	return nil
}

// GET https://www.bitstamp.net/api/v2/transactions/{currency_pair}/?time=day
func (c *HttpClient) V2Transactions(currencyPair string, timeParam string) (response []V2TransactionsResponse, err error) {
	if err = validateCurrencyPair(currencyPair); err != nil {
		return
	}
	urlPath := fmt.Sprintf("/v2/transactions/%s/", currencyPair)

	// quick n' dirty validation - from API docs:
	// The time interval from which we want the transactions to be returned. Possible values are minute, hour (default) or day.
	switch timeParam {
	case "":
		err = c.getRequest(&response, urlPath)
	case "minute", "hour", "day":
		err = c.getRequest(&response, urlPath, [2]string{"time", timeParam})
	default:
		err = fmt.Errorf("invalid value for time interval: %s", timeParam)
	}

	return
}

//
// Trading pairs info
//

type V2TradingPairsInfoResponse struct {
	BaseDecimals           int    `json:"base_decimals"`
	CounterDecimals        int    `json:"counter_decimals"`
	Description            string `json:"description"`
	InstantAndMarketOrders string `json:"instant_and_market_orders"` // TODO: make this a boolean flag
	MinimumOrder           string `json:"minimum_order"`
	Name                   string `json:"name"`
	Trading                string `json:"trading"` // TODO: make this a boolean flag
	UrlSymbol              string `json:"url_symbol"`
}

func (c *HttpClient) V2TradingPairsInfo() (response []V2TradingPairsInfoResponse, err error) {
	err = c.getRequest(&response, "/v2/trading-pairs-info/")
	return
}

//
// OHLC data
//

type Ohlc struct {
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Timestamp int64           `json:"timestamp"`
	Volume    decimal.Decimal `json:"volume"`
}

// custom deserialization instructions necessary
func (e *Ohlc) UnmarshalJSON(bytes []byte) error {
	if e == nil {
		e = new(Ohlc)
	}

	var raw map[string]string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	for _, k := range []string{"Open", "High", "Low", "Close", "Volume"} {
		kl := strings.ToLower(k)
		if strVal, exists := raw[kl]; exists {
			val, err2 := decimal.NewFromString(strVal)
			if err2 != nil {
				return fmt.Errorf("error parsing %s: %v", kl, err2)
			}

			tv := reflect.ValueOf(e).Elem().FieldByName(k)
			if tv.IsValid() {
				tv.Set(reflect.ValueOf(val))
			}
		}
	}

	if val, exists := raw["timestamp"]; exists {
		tstmp, err2 := strconv.ParseInt(val, 10, 64)
		if err2 != nil {
			return fmt.Errorf("error parsing timestamp: %v", err2)
		}
		e.Timestamp = tstmp
	} else {
		return fmt.Errorf("missing expected `date` field in %v", raw)
	}

	return nil
}

type V2OhlcResponse struct {
	Data struct {
		Pair    string `json:"pair"`
		Candles []Ohlc `json:"ohlc"`
	} `json:"data"`
}

// GET https://www.bitstamp.net/api/v2/ohlc/{currency_pair}/?step=60&limit=5
// - start (Optional): Unix timestamp from when OHLC data will be started.
// - end (Optional): Unix timestamp to when OHLC data will be shown.
// 	If none from start or end timestamps are posted then endpoint returns OHLC data to current unixtime. If both start and end timestamps are posted, end timestamp will be used.
// - step: Timeframe in seconds. Possible options are 60, 180, 300, 900, 1800, 3600, 7200, 14400, 21600, 43200, 86400, 259200
// - limit: Limit OHLC results (minimum: 1; maximum: 1000)
func (c *HttpClient) V2Ohlc(currencyPair string, step, limit int, start, end int64) (response V2OhlcResponse, err error) {
	if err = validateCurrencyPair(currencyPair); err != nil {
		return
	}

	validSteps := map[int]struct{}{
		60:     {},
		180:    {},
		300:    {},
		900:    {},
		1800:   {},
		3600:   {},
		7200:   {},
		14400:  {},
		21600:  {},
		43200:  {},
		86400:  {},
		259200: {},
	}

	args := make([][2]string, 0)

	if _, exists := validSteps[step]; exists {
		args = append(args, [2]string{"step", strconv.Itoa(step)})
	} else {
		err = fmt.Errorf("invalid value for step parameter: %d", step)
		return
	}

	if limit < 1 || limit > 1000 {
		err = fmt.Errorf("invalid value for limit parameter: %d", limit)
	} else {
		args = append(args, [2]string{"limit", strconv.Itoa(limit)})
	}

	if end != 0 {
		args = append(args, [2]string{"end", fmt.Sprintf("%d", end)})
	} else {
		if start != 0 {
			args = append(args, [2]string{"start", fmt.Sprintf("%d", start)})
		}
	}
	urlPath := fmt.Sprintf("/v2/ohlc/%s/", currencyPair)
	err = c.getRequest(&response, urlPath, args...)
	return
}

//
// EUR/USD conversion rate
//

type V2EurUsdResponse struct {
	Buy  decimal.Decimal `json:"buy"`
	Sell decimal.Decimal `json:"sell"`
}

// GET https://www.bitstamp.net/api/v2/eur_usd/
func (c *HttpClient) V2EurUsd() (response V2EurUsdResponse, err error) {
	err = c.getRequest(&response, "/v2/eur_usd/")
	return
}
