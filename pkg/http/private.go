package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
)

// Contains "private" endpoints whereby we are following the naming here: https://www.bitstamp.net/api/

//
// Account balance
//

type V2BalanceResponse struct {
	// currencies
	AaveAvailable      decimal.Decimal `json:"aave_available"`
	AaveBalance        decimal.Decimal `json:"aave_balance"`
	AaveReserved       decimal.Decimal `json:"aave_reserved"`
	AaveWithdrawalFee  decimal.Decimal `json:"aave_withdrawal_fee"`
	AlgoAvailable      decimal.Decimal `json:"algo_available"`
	AlgoBalance        decimal.Decimal `json:"algo_balance"`
	AlgoReserved       decimal.Decimal `json:"algo_reserved"`
	AlgoWithdrawalFee  decimal.Decimal `json:"algo_withdrawal_fee"`
	AudioAvailable     decimal.Decimal `json:"audio_available"`
	AudioBalance       decimal.Decimal `json:"audio_balance"`
	AudioReserved      decimal.Decimal `json:"audio_reserved"`
	AudioWithdrawalFee decimal.Decimal `json:"audio_withdrawal_fee"`
	BatAvailable       decimal.Decimal `json:"bat_available"`
	BatBalance         decimal.Decimal `json:"bat_balance"`
	BatReserved        decimal.Decimal `json:"bat_reserved"`
	BatWithdrawalFee   decimal.Decimal `json:"bat_withdrawal_fee"`
	BchAvailable       decimal.Decimal `json:"bch_available"`
	BchBalance         decimal.Decimal `json:"bch_balance"`
	BchReserved        decimal.Decimal `json:"bch_reserved"`
	BchWithdrawalFee   decimal.Decimal `json:"bch_withdrawal_fee"`
	BtcAvailable       decimal.Decimal `json:"btc_available"`
	BtcBalance         decimal.Decimal `json:"btc_balance"`
	BtcReserved        decimal.Decimal `json:"btc_reserved"`
	BtcWithdrawalFee   decimal.Decimal `json:"btc_withdrawal_fee"`
	CompAvailable      decimal.Decimal `json:"comp_available"`
	CompBalance        decimal.Decimal `json:"comp_balance"`
	CompReserved       decimal.Decimal `json:"comp_reserved"`
	CompWithdrawalFee  decimal.Decimal `json:"comp_withdrawal_fee"`
	CrvAvailable       decimal.Decimal `json:"crv_available"`
	CrvBalance         decimal.Decimal `json:"crv_balance"`
	CrvReserved        decimal.Decimal `json:"crv_reserved"`
	CrvWithdrawalFee   decimal.Decimal `json:"crv_withdrawal_fee"`
	DaiAvailable       decimal.Decimal `json:"dai_available"`
	DaiBalance         decimal.Decimal `json:"dai_balance"`
	DaiReserved        decimal.Decimal `json:"dai_reserved"`
	DaiWithdrawalFee   decimal.Decimal `json:"dai_withdrawal_fee"`
	Eth2Available      decimal.Decimal `json:"eth2_available"`
	Eth2Balance        decimal.Decimal `json:"eth2_balance"`
	Eth2Reserved       decimal.Decimal `json:"eth2_reserved"`
	Eth2rAvailable     decimal.Decimal `json:"eth2r_available"`
	Eth2rBalance       decimal.Decimal `json:"eth2r_balance"`
	Eth2rReserved      decimal.Decimal `json:"eth2r_reserved"`
	EthAvailable       decimal.Decimal `json:"eth_available"`
	EthBalance         decimal.Decimal `json:"eth_balance"`
	EthReserved        decimal.Decimal `json:"eth_reserved"`
	EthWithdrawalFee   decimal.Decimal `json:"eth_withdrawal_fee"`
	EurAvailable       decimal.Decimal `json:"eur_available"`
	EurBalance         decimal.Decimal `json:"eur_balance"`
	EurReserved        decimal.Decimal `json:"eur_reserved"`
	EurWithdrawalFee   decimal.Decimal `json:"eur_withdrawal_fee"`
	GbpAvailable       decimal.Decimal `json:"gbp_available"`
	GbpBalance         decimal.Decimal `json:"gbp_balance"`
	GbpReserved        decimal.Decimal `json:"gbp_reserved"`
	GbpWithdrawalFee   decimal.Decimal `json:"gbp_withdrawal_fee"`
	GrtAvailable       decimal.Decimal `json:"grt_available"`
	GrtBalance         decimal.Decimal `json:"grt_balance"`
	GrtReserved        decimal.Decimal `json:"grt_reserved"`
	GrtWithdrawalFee   decimal.Decimal `json:"grt_withdrawal_fee"`
	GusdAvailable      decimal.Decimal `json:"gusd_available"`
	GusdBalance        decimal.Decimal `json:"gusd_balance"`
	GusdReserved       decimal.Decimal `json:"gusd_reserved"`
	GusdWithdrawalFee  decimal.Decimal `json:"gusd_withdrawal_fee"`
	KncAvailable       decimal.Decimal `json:"knc_available"`
	KncBalance         decimal.Decimal `json:"knc_balance"`
	KncReserved        decimal.Decimal `json:"knc_reserved"`
	KncWithdrawalFee   decimal.Decimal `json:"knc_withdrawal_fee"`
	LinkAvailable      decimal.Decimal `json:"link_available"`
	LinkBalance        decimal.Decimal `json:"link_balance"`
	LinkReserved       decimal.Decimal `json:"link_reserved"`
	LinkWithdrawalFee  decimal.Decimal `json:"link_withdrawal_fee"`
	LtcAvailable       decimal.Decimal `json:"ltc_available"`
	LtcBalance         decimal.Decimal `json:"ltc_balance"`
	LtcReserved        decimal.Decimal `json:"ltc_reserved"`
	LtcWithdrawalFee   decimal.Decimal `json:"ltc_withdrawal_fee"`
	MkrAvailable       decimal.Decimal `json:"mkr_available"`
	MkrBalance         decimal.Decimal `json:"mkr_balance"`
	MkrReserved        decimal.Decimal `json:"mkr_reserved"`
	MkrWithdrawalFee   decimal.Decimal `json:"mkr_withdrawal_fee"`
	OmgAvailable       decimal.Decimal `json:"omg_available"`
	OmgBalance         decimal.Decimal `json:"omg_balance"`
	OmgReserved        decimal.Decimal `json:"omg_reserved"`
	OmgWithdrawalFee   decimal.Decimal `json:"omg_withdrawal_fee"`
	PaxAvailable       decimal.Decimal `json:"pax_available"`
	PaxBalance         decimal.Decimal `json:"pax_balance"`
	PaxReserved        decimal.Decimal `json:"pax_reserved"`
	PaxWithdrawalFee   decimal.Decimal `json:"pax_withdrawal_fee"`
	SnxAvailable       decimal.Decimal `json:"snx_available"`
	SnxBalance         decimal.Decimal `json:"snx_balance"`
	SnxReserved        decimal.Decimal `json:"snx_reserved"`
	SnxWithdrawalFee   decimal.Decimal `json:"snx_withdrawal_fee"`
	UmaAvailable       decimal.Decimal `json:"uma_available"`
	UmaBalance         decimal.Decimal `json:"uma_balance"`
	UmaReserved        decimal.Decimal `json:"uma_reserved"`
	UmaWithdrawalFee   decimal.Decimal `json:"uma_withdrawal_fee"`
	UniAvailable       decimal.Decimal `json:"uni_available"`
	UniBalance         decimal.Decimal `json:"uni_balance"`
	UniReserved        decimal.Decimal `json:"uni_reserved"`
	UniWithdrawalFee   decimal.Decimal `json:"uni_withdrawal_fee"`
	UsdAvailable       decimal.Decimal `json:"usd_available"`
	UsdBalance         decimal.Decimal `json:"usd_balance"`
	UsdReserved        decimal.Decimal `json:"usd_reserved"`
	UsdWithdrawalFee   decimal.Decimal `json:"usd_withdrawal_fee"`
	UsdcAvailable      decimal.Decimal `json:"usdc_available"`
	UsdcBalance        decimal.Decimal `json:"usdc_balance"`
	UsdcReserved       decimal.Decimal `json:"usdc_reserved"`
	UsdcWithdrawalFee  decimal.Decimal `json:"usdc_withdrawal_fee"`
	UsdtAvailable      decimal.Decimal `json:"usdt_available"`
	UsdtBalance        decimal.Decimal `json:"usdt_balance"`
	UsdtReserved       decimal.Decimal `json:"usdt_reserved"`
	UsdtWithdrawalFee  decimal.Decimal `json:"usdt_withdrawal_fee"`
	XlmAvailable       decimal.Decimal `json:"xlm_available"`
	XlmBalance         decimal.Decimal `json:"xlm_balance"`
	XlmReserved        decimal.Decimal `json:"xlm_reserved"`
	XlmWithdrawalFee   decimal.Decimal `json:"xlm_withdrawal_fee"`
	XrpAvailable       decimal.Decimal `json:"xrp_available"`
	XrpBalance         decimal.Decimal `json:"xrp_balance"`
	XrpReserved        decimal.Decimal `json:"xrp_reserved"`
	XrpWithdrawalFee   decimal.Decimal `json:"xrp_withdrawal_fee"`
	YfiAvailable       decimal.Decimal `json:"yfi_available"`
	YfiBalance         decimal.Decimal `json:"yfi_balance"`
	YfiReserved        decimal.Decimal `json:"yfi_reserved"`
	YfiWithdrawalFee   decimal.Decimal `json:"yfi_withdrawal_fee"`
	ZrxAvailable       decimal.Decimal `json:"zrx_available"`
	ZrxBalance         decimal.Decimal `json:"zrx_balance"`
	ZrxReserved        decimal.Decimal `json:"zrx_reserved"`
	ZrxWithdrawalFee   decimal.Decimal `json:"zrx_withdrawal_fee"`

	// pairs
	AavebtcFee  decimal.Decimal `json:"aavebtc_fee"`
	AaveeurFee  decimal.Decimal `json:"aaveeur_fee"`
	AaveusdFee  decimal.Decimal `json:"aaveusd_fee"`
	AlgobtcFee  decimal.Decimal `json:"algobtc_fee"`
	AlgoeurFee  decimal.Decimal `json:"algoeur_fee"`
	AlgousdFee  decimal.Decimal `json:"algousd_fee"`
	AudiobtcFee decimal.Decimal `json:"audiobtc_fee"`
	AudioeurFee decimal.Decimal `json:"audioeur_fee"`
	AudiousdFee decimal.Decimal `json:"audiousd_fee"`
	BatbtcFee   decimal.Decimal `json:"batbtc_fee"`
	BateurFee   decimal.Decimal `json:"bateur_fee"`
	BatusdFee   decimal.Decimal `json:"batusd_fee"`
	BchbtcFee   decimal.Decimal `json:"bchbtc_fee"`
	BcheurFee   decimal.Decimal `json:"bcheur_fee"`
	BchgbpFee   decimal.Decimal `json:"bchgbp_fee"`
	BchusdFee   decimal.Decimal `json:"bchusd_fee"`
	BtceurFee   decimal.Decimal `json:"btceur_fee"`
	BtcgbpFee   decimal.Decimal `json:"btcgbp_fee"`
	BtcpaxFee   decimal.Decimal `json:"btcpax_fee"`
	BtcusdFee   decimal.Decimal `json:"btcusd_fee"`
	BtcusdcFee  decimal.Decimal `json:"btcusdc_fee"`
	BtcusdtFee  decimal.Decimal `json:"btcusdt_fee"`
	CompbtcFee  decimal.Decimal `json:"compbtc_fee"`
	CompeurFee  decimal.Decimal `json:"compeur_fee"`
	CompusdFee  decimal.Decimal `json:"compusd_fee"`
	CrvbtcFee   decimal.Decimal `json:"crvbtc_fee"`
	CrveurFee   decimal.Decimal `json:"crveur_fee"`
	CrvusdFee   decimal.Decimal `json:"crvusd_fee"`
	DaiusdFee   decimal.Decimal `json:"daiusd_fee"`
	Eth2ethFee  decimal.Decimal `json:"eth2eth_fee"`
	EthbtcFee   decimal.Decimal `json:"ethbtc_fee"`
	EtheurFee   decimal.Decimal `json:"etheur_fee"`
	EthgbpFee   decimal.Decimal `json:"ethgbp_fee"`
	EthpaxFee   decimal.Decimal `json:"ethpax_fee"`
	EthusdFee   decimal.Decimal `json:"ethusd_fee"`
	EthusdcFee  decimal.Decimal `json:"ethusdc_fee"`
	EthusdtFee  decimal.Decimal `json:"ethusdt_fee"`
	EurusdFee   decimal.Decimal `json:"eurusd_fee"`
	GbpeurFee   decimal.Decimal `json:"gbpeur_fee"`
	GbpusdFee   decimal.Decimal `json:"gbpusd_fee"`
	GusdusdFee  decimal.Decimal `json:"gusdusd_fee"`
	KncbtcFee   decimal.Decimal `json:"kncbtc_fee"`
	KnceurFee   decimal.Decimal `json:"knceur_fee"`
	KncusdFee   decimal.Decimal `json:"kncusd_fee"`
	LinkbtcFee  decimal.Decimal `json:"linkbtc_fee"`
	LinkethFee  decimal.Decimal `json:"linketh_fee"`
	LinkeurFee  decimal.Decimal `json:"linkeur_fee"`
	LinkgbpFee  decimal.Decimal `json:"linkgbp_fee"`
	LinkusdFee  decimal.Decimal `json:"linkusd_fee"`
	LtcbtcFee   decimal.Decimal `json:"ltcbtc_fee"`
	LtceurFee   decimal.Decimal `json:"ltceur_fee"`
	LtcgbpFee   decimal.Decimal `json:"ltcgbp_fee"`
	LtcusdFee   decimal.Decimal `json:"ltcusd_fee"`
	MkrbtcFee   decimal.Decimal `json:"mkrbtc_fee"`
	MkreurFee   decimal.Decimal `json:"mkreur_fee"`
	MkrusdFee   decimal.Decimal `json:"mkrusd_fee"`
	OmgbtcFee   decimal.Decimal `json:"omgbtc_fee"`
	OmgeurFee   decimal.Decimal `json:"omgeur_fee"`
	OmggbpFee   decimal.Decimal `json:"omggbp_fee"`
	OmgusdFee   decimal.Decimal `json:"omgusd_fee"`
	PaxeurFee   decimal.Decimal `json:"paxeur_fee"`
	PaxgbpFee   decimal.Decimal `json:"paxgbp_fee"`
	PaxusdFee   decimal.Decimal `json:"paxusd_fee"`
	SnxbtcFee   decimal.Decimal `json:"snxbtc_fee"`
	SnxeurFee   decimal.Decimal `json:"snxeur_fee"`
	SnxusdFee   decimal.Decimal `json:"snxusd_fee"`
	UmabtcFee   decimal.Decimal `json:"umabtc_fee"`
	UmaeurFee   decimal.Decimal `json:"umaeur_fee"`
	UmausdFee   decimal.Decimal `json:"umausd_fee"`
	UnibtcFee   decimal.Decimal `json:"unibtc_fee"`
	UnieurFee   decimal.Decimal `json:"unieur_fee"`
	UniusdFee   decimal.Decimal `json:"uniusd_fee"`
	UsdceurFee  decimal.Decimal `json:"usdceur_fee"`
	UsdcusdFee  decimal.Decimal `json:"usdcusd_fee"`
	UsdcusdtFee decimal.Decimal `json:"usdcusdt_fee"`
	UsdteurFee  decimal.Decimal `json:"usdteur_fee"`
	UsdtusdFee  decimal.Decimal `json:"usdtusd_fee"`
	XlmbtcFee   decimal.Decimal `json:"xlmbtc_fee"`
	XlmeurFee   decimal.Decimal `json:"xlmeur_fee"`
	XlmgbpFee   decimal.Decimal `json:"xlmgbp_fee"`
	XlmusdFee   decimal.Decimal `json:"xlmusd_fee"`
	XrpbtcFee   decimal.Decimal `json:"xrpbtc_fee"`
	XrpeurFee   decimal.Decimal `json:"xrpeur_fee"`
	XrpgbpFee   decimal.Decimal `json:"xrpgbp_fee"`
	XrppaxFee   decimal.Decimal `json:"xrppax_fee"`
	XrpusdFee   decimal.Decimal `json:"xrpusd_fee"`
	XrpusdtFee  decimal.Decimal `json:"xrpusdt_fee"`
	YfibtcFee   decimal.Decimal `json:"yfibtc_fee"`
	YfieurFee   decimal.Decimal `json:"yfieur_fee"`
	YfiusdFee   decimal.Decimal `json:"yfiusd_fee"`
	ZrxbtcFee   decimal.Decimal `json:"zrxbtc_fee"`
	ZrxeurFee   decimal.Decimal `json:"zrxeur_fee"`
	ZrxusdFee   decimal.Decimal `json:"zrxusd_fee"`

	// fee
	Fee decimal.Decimal `json:"fee"`
}

// POST https://www.bitstamp.net/api/v2/balance/
// POST https://www.bitstamp.net/api/v2/balance/{currency_pair}/
func (c *HttpClient) V2Balance(currencyPairOrAll string) (response V2BalanceResponse, err error) {
	// TODO: validate currency pair
	if currencyPairOrAll == "all" {
		err = c.authenticatedPostRequest(&response, "/v2/balance/")
	} else {
		err = c.authenticatedPostRequest(&response, fmt.Sprintf("/v2/balance/%s/", currencyPairOrAll))
	}

	return
}

// User transactions
type V2UserTransactionsResponse struct {
	Datetime string          `json:"datetime"`
	Fee      decimal.Decimal `json:"fee"`
	Id       int64           `json:"id"`
	OrderId  int64           `json:"order_id"`
	Type     string          `json:"type"`

	Status string      `json:"status"`
	Reason interface{} `json:"reason"`

	// amounts
	Aave   decimal.Decimal `json:"aave"`
	Algo   decimal.Decimal `json:"algo"`
	Audio  decimal.Decimal `json:"audio"`
	Bat    decimal.Decimal `json:"bat"`
	Bch    decimal.Decimal `json:"bch"`
	Btc    decimal.Decimal `json:"btc"`
	Comp   decimal.Decimal `json:"comp"`
	Crv    decimal.Decimal `json:"crv"`
	Dai    decimal.Decimal `json:"dai"`
	Eth    decimal.Decimal `json:"eth"`
	Eth2   decimal.Decimal `json:"eth2"`
	Eth2r  decimal.Decimal `json:"eth2r"`
	Eur    decimal.Decimal `json:"eur"`
	Gbp    decimal.Decimal `json:"gbp"`
	Grt    decimal.Decimal `json:"grt"`
	Gusd   decimal.Decimal `json:"gusd"`
	Knc    decimal.Decimal `json:"knc"`
	Link   decimal.Decimal `json:"link"`
	Ltc    decimal.Decimal `json:"ltc"`
	Mkr    decimal.Decimal `json:"mkr"`
	Omg    decimal.Decimal `json:"omg"`
	Pax    decimal.Decimal `json:"pax"`
	Snx    decimal.Decimal `json:"snx"`
	Uma    decimal.Decimal `json:"uma"`
	Uni    decimal.Decimal `json:"uni"`
	Usd    decimal.Decimal `json:"usd"`
	Usdc   decimal.Decimal `json:"usdc"`
	Usdt   decimal.Decimal `json:"usdt"`
	Xlm    decimal.Decimal `json:"xlm"`
	Xrp    decimal.Decimal `json:"xrp"`
	Yfi    decimal.Decimal `json:"yfi"`
	Zrx    decimal.Decimal `json:"zrx"`
	BtcUsd decimal.Decimal `json:"btc_usd"`
}

// TODO: add arguments!
func (c *HttpClient) V2UserTransactions(currencyPairOrAll string) (response []V2UserTransactionsResponse, err error) {
	if currencyPairOrAll == "all" {
		err = c.authenticatedPostRequest(&response, "/v2/user_transactions/", [2]string{"limit", "1000"})
	} else {
		err = c.authenticatedPostRequest(&response, fmt.Sprintf("/v2/user_transactions/%s/", currencyPairOrAll), [2]string{"limit", "1000"})
	}

	return
}

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
func (c *HttpClient) V2OpenOrders(currencyPairOrAll string) (response []V2OpenOrdersResponse, err error) {
	urlPath := fmt.Sprintf("/v2/open_orders/%s/", currencyPairOrAll)
	err = c.authenticatedPostRequest(&response, urlPath)

	return
}

//
// Order status
//

type V2OrderStatusTransaction struct {
	Btc      decimal.Decimal `json:"btc"`
	Datetime string          `json:"datetime"`
	Fee      decimal.Decimal `json:"fee"`
	Price    decimal.Decimal `json:"price"`
	Tid      int64           `json:"tid"`
	Type     int             `json:"type"`
	Usd      decimal.Decimal `json:"usd"`
}

type V2OrderStatusResponse struct {
	Id              int64                      `json:"id"`
	Status          string                     `json:"status"`
	Transactions    []V2OrderStatusTransaction `json:"transactions"`
	AmountRemaining decimal.Decimal            `json:"amount_remaining"`
	ClientOrderId   string                     `json:"client_order_id"`

	Reason interface{} `json:"reason"`
}

// POST https://www.bitstamp.net/api/v2/order_status/
func (c *HttpClient) V2OrderStatus(orderId int64, clOrdId string, omitTx bool) (response V2OrderStatusResponse, err error) {
	params := make([][2]string, 0)
	params = append(params, [2]string{"id", fmt.Sprintf("%d", orderId)})
	if clOrdId != "" {
		params = append(params, [2]string{"client_order_id", clOrdId})
	}
	if omitTx {
		params = append(params, [2]string{"omit_transactions", "true"})
	}

	err = c.authenticatedPostRequest(&response, "/v2/order_status/", params...)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error: %v", response.Reason)
	}

	return
}

//
// Cancel order
//

type V2CancelOrderResponse struct {
	Id     uint64          `json:"id"`
	Amount decimal.Decimal `json:"amount"`
	Price  decimal.Decimal `json:"price"`
	Type   uint8           `json:"type"`
	Error  string          `json:"error"`
}

func (c *HttpClient) V2CancelOrder(orderId int64) (response V2CancelOrderResponse, err error) {
	err = c.authenticatedPostRequest(&response, "/v2/cancel_order/", [2]string{"id", fmt.Sprintf("%d", orderId)})
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

func (c *HttpClient) v2LimitOrder(side, currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s/%s/", side, currencyPair)

	if c.autoRounding {
		// TODO: we probably need "smarter" (stingier?) rounding here...
		amount = amount.Round(roundings[currencyPair].Base)
		price = price.Round(roundings[currencyPair].Counter)
	}

	params := make([][2]string, 0)
	params = append(params, [2]string{"amount", amount.String()})
	params = append(params, [2]string{"price", price.String()})
	if dailyOrder {
		params = append(params, [2]string{"daily_order", "True"})
	}
	if iocOrder {
		params = append(params, [2]string{"ioc_order", "True"})
	}
	if clOrdId != "" {
		params = append(params, [2]string{"client_order_id", clOrdId})
	}
	// TODO: limitPrice !

	err = c.authenticatedPostRequest(&response, urlPath, params...)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error placing limit %s (%s @ %s): %v", side, amount, price, response.Reason)
	}

	return
}

func (c *HttpClient) V2BuyLimitOrder(currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
	return c.v2LimitOrder("buy", currencyPair, price, amount, limitPrice, dailyOrder, iocOrder, clOrdId)
}

func (c *HttpClient) V2SellLimitOrder(currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
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

func (c *HttpClient) v2MarketOrder(side, currencyPair string, amount decimal.Decimal, clOrdId string) (response V2MarketOrderResponse, err error) {
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

func (c *HttpClient) V2BuyMarketOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2MarketOrderResponse, err error) {
	return c.v2MarketOrder("buy", currencyPair, amount, clOrdId)
}

func (c *HttpClient) V2SellMarketOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2MarketOrderResponse, err error) {
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

func (c *HttpClient) v2InstantOrder(side, currencyPair string, amount decimal.Decimal, clOrdId string) (response V2InstantOrderResponse, err error) {
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

func (c *HttpClient) V2BuyInstantOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2InstantOrderResponse, err error) {
	return c.v2InstantOrder("buy", currencyPair, amount, clOrdId)
}

func (c *HttpClient) V2SellInstantOrder(currencyPair string, amount decimal.Decimal, clOrdId string) (response V2InstantOrderResponse, err error) {
	return c.v2InstantOrder("sell", currencyPair, amount, clOrdId)
}

type V2WebsocketsTokenResponse struct {
	Token    string `json:"token"`
	ValidSec uint32 `json:"valid_sec"`
	UserId   uint32 `json:"user_id"`
}

// V2WebsocketsToken generates an ephemeral token, which allows user to subscribe to private
// websocket events. These events include ClientOrderIds (and potentially additional private data)
func (c *HttpClient) V2WebsocketsToken() (response V2WebsocketsTokenResponse, err error) {
	err = c.authenticatedPostRequest(&response, "/v2/websockets_token/")
	if err != nil {
		return
	}
	return
}
