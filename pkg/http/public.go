package http

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
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
	urlPath := fmt.Sprintf("/v2/ticker/%s/", currencyPair)
	err = c.getRequest(&response, urlPath)
	return
}

// GET https://www.bitstamp.net/api/v2/ticker_hour/{currency_pair}/
func (c *HttpClient) V2HourlyTicker(currencyPair string) (response TickerResponse, err error) {
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
// GET https://www.bitstamp.net/api/order_book?group=1
func (c *HttpClient) V2OrderBook(currencyPair string, group int) (response V2OrderBookResponse, err error) {
	urlPath := fmt.Sprintf("/v2/order_book/%s/", currencyPair)
	err = c.getRequest(&response, urlPath, [2]string{"group", strconv.Itoa(group)})
	return
}

//
// Transactions
//

//
// Trading pairs info
//

//
// OHLC data
//

//
// EUR/USD conversion rate
//
