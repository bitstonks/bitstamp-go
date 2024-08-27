package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// StringInt create a type alias for type int64
type StringInt int64

// UnmarshalJSON create a custom unmarshal for the StringInt
// / this helps us check the type of our value before unmarshalling it
func (st *StringInt) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into a int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int64(v))
	case string:
		///here convert the string into
		///an integer
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			///the string might not be of integer type
			///so return an error
			return err

		}
		*st = StringInt(i)

	}
	return nil
}

type StringDatetimeFromMicroseconds time.Time

func (st *StringDatetimeFromMicroseconds) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringDatetimeFromMicroseconds(time.Unix(0, int64(v)*int64(time.Microsecond)))
	case float64:
		*st = StringDatetimeFromMicroseconds(time.Unix(0, int64(v)*int64(time.Microsecond)))
	case string:
		millis, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			///the string might not be of integer type
			///so return an error
			return err

		}
		*st = StringDatetimeFromMicroseconds(time.Unix(0, millis*int64(time.Microsecond)))

	}
	return nil
}

// Contains "private" endpoints whereby we are following the naming here: https://www.bitstamp.net/api/

// Balance

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
		err = c.authenticatedFormRequest(&response, "POST", "/v2/balance/", nil, nil)
	} else {
		err = c.authenticatedFormRequest(&response, "POST", fmt.Sprintf("/v2/balance/%s/", currencyPairOrAll), nil, nil)
	}

	return
}

// Account Balances

type V2AccountBalancesResponse struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Reserved  decimal.Decimal `json:"reserved"`
	Total     decimal.Decimal `json:"total"`
}

// POST https://www.bitstamp.net/api/v2/account_balances/
func (c *HttpClient) V2AccountBalances() (response []V2AccountBalancesResponse, err error) {
	err = c.authenticatedFormRequest(&response, "POST", "/v2/account_balances/", nil, nil)
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
		err = c.authenticatedFormRequest(&response, "POST", "/v2/user_transactions/", nil, map[string]string{"limit": "1000"})
	} else {
		err = c.authenticatedFormRequest(&response, "POST", fmt.Sprintf("/v2/user_transactions/%s/", currencyPairOrAll), nil, map[string]string{"limit": "1000"})
	}

	return
}

// Crypto Transactions

type V2CryptoTransactionDepositWithdrawal struct {
	Datetime           string          `json:"datetime"`
	Txid               string          `json:"txid"`
	DesitnationAddress string          `json:"desitnationAddress"`
	Amount             decimal.Decimal `json:"amount"`
	Network            string          `json:"network"`
	Currency           string          `json:"currency"`
}

type V2CryptoTransactionIou struct {
	Datetime           string          `json:"datetime"`
	Txid               string          `json:"txid"`
	DesitnationAddress string          `json:"desitnationAddress"`
	Amount             decimal.Decimal `json:"amount"`
	Network            string          `json:"network"`
	Currency           string          `json:"currency"`
	Type               string          `json:"type"`
}

type V2CryptoTransactionsResponse struct {
	Deposits              []V2CryptoTransactionDepositWithdrawal `json:"deposits"`
	Withdrawals           []V2CryptoTransactionDepositWithdrawal `json:"withdrawals"`
	RippleIouTransactions []V2CryptoTransactionIou               `json:"ripple_iou_transactions"`
	Status                string                                 `json:"status"`
	Reason                interface{}                            `json:"reason"`
}

func (c *HttpClient) V2CryptoTransactions(includeIous bool) (response V2CryptoTransactionsResponse, err error) {
	params := map[string]string{"limit": "1000"}
	if includeIous {
		params["include_ious"] = ""
	}

	err = c.authenticatedFormRequest(&response, "POST", "/v2/crypto-transactions/", nil, params)
	return
}

// Crypto Address

type V2CryptoAddressResponse struct {
	Address string `json:"address"`
	Error   string `json:"error"`
}

func (c *HttpClient) V2CryptoAddress(currency string) (response V2CryptoAddressResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s_address/", currency)
	err = c.authenticatedFormRequest(&response, "POST", urlPath, nil, nil)
	return
}

// Withdrawal Requests

type V2WithdrawalRequestsResponse struct {
	Id       string          `json:"id"`
	Datetime string          `json:"datetime"`
	Type     string          `json:"type"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
	Status   string          `json:"status"`
	Txid     string          `json:"txid"`
	Reason   interface{}     `json:"reason"`
}

func (c *HttpClient) V2WithdrawalRequests(withdrawalId int64, timeDelta string) (response []V2WithdrawalRequestsResponse, err error) {
	params := map[string]string{"offset": "0", "limit": "1000"}
	if withdrawalId != 0 {
		params["id"] = fmt.Sprintf("%d", withdrawalId)
	}
	if timeDelta != "" {
		params["timedelta"] = ""
	}

	err = c.authenticatedFormRequest(&response, "POST", "/v2/withdrawal-requests/", nil, params)
	return
}

// Withdrawal Fees

type V2WithdrawalFeesResponse struct {
	Currency string          `json:"currency"`
	Fee      decimal.Decimal `json:"amount"`
	Network  string          `json:"status"`
}

func (c *HttpClient) V2WithdrawalFees() (response []V2WithdrawalFeesResponse, err error) {
	err = c.authenticatedFormRequest(&response, "POST", "/v2/fees/withdrawal/", nil, nil)
	return
}

// Trading Fees

type V2TradingFees struct {
	Maker decimal.Decimal `json:"maker"`
	Taker decimal.Decimal `json:"taker"`
}

type V2TradingFeesResponse struct {
	CurrencyPair string        `json:"currency_pair"`
	Market       string        `json:"market"`
	Fees         V2TradingFees `json:"fees"`
}

func (c *HttpClient) V2TradingFees() (response []V2TradingFeesResponse, err error) {
	err = c.authenticatedFormRequest(&response, "POST", "/v2/fees/trading/", nil, nil)
	return
}

// Open orders
type V2OpenOrdersResponse struct {
	Id            string           `json:"id"`
	Datetime      string           `json:"datetime"`
	Type          string           `json:"type"`
	Price         decimal.Decimal  `json:"price"`
	Amount        decimal.Decimal  `json:"amount"`
	CurrencyPair  string           `json:"currency_pair"`
	ClientOrderId string           `json:"client_order_id"`
	Status        string           `json:"status"`
	Reason        interface{}      `json:"reason"`
	Leverage      *decimal.Decimal `json:"leverage"`
	MarginMode    *MarginMode      `json:"margin_mode"`
}

// POST https://www.bitstamp.net/api/v2/open_orders/all/
// POST https://www.bitstamp.net/api/v2/open_orders/{currency_pair}
func (c *HttpClient) V2OpenOrders(currencyPairOrAll string) (response []V2OpenOrdersResponse, err error) {
	urlPath := fmt.Sprintf("/v2/open_orders/%s/", currencyPairOrAll)
	err = c.authenticatedFormRequest(&response, "POST", urlPath, nil, nil)

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
	params := map[string]string{
		"id": fmt.Sprintf("%d", orderId),
	}
	if clOrdId != "" {
		params["client_order_id"] = clOrdId
	}
	if omitTx {
		params["omit_transactions"] = "true"
	}

	err = c.authenticatedFormRequest(&response, "POST", "/v2/order_status/", nil, params)
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
	err = c.authenticatedFormRequest(&response, "POST", "/v2/cancel_order/", nil, map[string]string{"id": fmt.Sprintf("%d", orderId)})
	return
}

// Cancel all orders
// Buy limit order
// Sell limit order

// {"status": "error", "reason": {"__all__": ["Price is more than 20% below market price."]}}
// {"status": "error", "reason": {"__all__": ["You need 158338.86 USD to open that order. You have only 99991.52 USD available. Check your account balance for details."]}}
type V2LimitOrderResponse struct {
	Id         string           `json:"id"`
	Datetime   string           `json:"datetime"`
	Type       string           `json:"type"`
	Price      decimal.Decimal  `json:"price"`
	Amount     decimal.Decimal  `json:"amount"`
	Status     string           `json:"status"`
	Reason     interface{}      `json:"reason"`
	Leverage   *decimal.Decimal `json:"leverage"`
	MarginMode *MarginMode      `json:"margin_mode"`
}

type MarginMode string

const (
	Cross    MarginMode = "CROSS"
	Isolated MarginMode = "ISOLATED"
)

func (c *HttpClient) v2LimitOrder(side, currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2LimitOrderResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s/%s/", side, currencyPair)

	if c.autoRounding {
		// TODO: we probably need "smarter" (stingier?) rounding here...
		amount = amount.Round(roundings[currencyPair].Base)
		price = price.Round(roundings[currencyPair].Counter)
	}

	params := map[string]string{
		"amount": amount.String(),
		"price":  price.String(),
	}

	if dailyOrder {
		params["daily_order"] = "True"
	}
	if iocOrder {
		params["ioc_order"] = "True"
	}
	if clOrdId != "" {
		params["client_order_id"] = clOrdId
	}
	if marginMode != nil {
		params["margin_mode"] = string(*marginMode)
	}
	if leverage != nil {
		params["leverage"] = leverage.String()
	}
	if reduceOnly {
		params["reduce_only"] = "True"
	}
	// TODO: limitPrice !

	err = c.authenticatedFormRequest(&response, "POST", urlPath, nil, params)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error placing limit %s (%s @ %s): %v", side, amount, price, response.Reason)
	}

	return
}

func (c *HttpClient) V2BuyLimitOrder(currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2LimitOrderResponse, err error) {
	return c.v2LimitOrder("buy", currencyPair, price, amount, limitPrice, dailyOrder, iocOrder, clOrdId, marginMode, leverage, reduceOnly)
}

func (c *HttpClient) V2SellLimitOrder(currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2LimitOrderResponse, err error) {
	return c.v2LimitOrder("sell", currencyPair, price, amount, limitPrice, dailyOrder, iocOrder, clOrdId, marginMode, leverage, reduceOnly)
}

type V2MarketOrderResponse struct {
	Id              string          `json:"id"`
	Subtype         string          `json:"subtype"`
	Market          string          `json:"market"`
	Datetime        string          `json:"datetime"`
	Type            string          `json:"type"`
	Price           decimal.Decimal `json:"price"`
	Amount          decimal.Decimal `json:"amount"`
	ClientOrderId   string          `json:"client_order_id"`
	MarginMode      MarginMode      `json:"margin_mode"`
	Leverage        decimal.Decimal `json:"leverage"`
	StopPrice       decimal.Decimal `json:"stop_price"`
	Trigger         string          `json:"trigger"`
	ActivationPrice decimal.Decimal `json:"activation_price"`
	TrailingDelta   decimal.Decimal `json:"trailing_delta"`
	Reason          string          `json:"reason"`
	Status          string          `json:"status"`
}

func (c *HttpClient) v2MarketOrder(side, currencyPair string, amount decimal.Decimal, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2MarketOrderResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s/market/%s/", side, currencyPair)

	data := make(map[string]string)
	data["amount"] = amount.String()
	if clOrdId != "" {
		data["client_order_id"] = clOrdId
	}
	if marginMode != nil {
		data["margin_mode"] = string(*marginMode)
	}
	if leverage != nil {
		data["leverage"] = leverage.String()
	}
	if reduceOnly {
		data["reduce_only"] = "True"
	}

	err = c.authenticatedFormRequest(response, "POST", urlPath, nil, data)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error placing market %s (for %s): %v", side, amount, response.Reason)
		return
	}

	return
}

func (c *HttpClient) V2BuyMarketOrder(currencyPair string, amount decimal.Decimal, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2MarketOrderResponse, err error) {
	return c.v2MarketOrder("buy", currencyPair, amount, clOrdId, marginMode, leverage, reduceOnly)
}

func (c *HttpClient) V2SellMarketOrder(currencyPair string, amount decimal.Decimal, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2MarketOrderResponse, err error) {
	return c.v2MarketOrder("sell", currencyPair, amount, clOrdId, marginMode, leverage, reduceOnly)
}

type V2InstantOrderResponse struct {
	Id         string           `json:"id"`
	Datetime   string           `json:"datetime"`
	Type       string           `json:"type"`
	Price      decimal.Decimal  `json:"price"`
	Amount     decimal.Decimal  `json:"amount"`
	Error      string           `json:"error"`
	Status     string           `json:"status"`
	Reason     interface{}      `json:"reason"`
	Leverage   *decimal.Decimal `json:"leverage"`
	MarginMode *MarginMode      `json:"margin_mode"`
}

func (c *HttpClient) v2InstantOrder(side, currencyPair string, amount decimal.Decimal, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2InstantOrderResponse, err error) {
	urlPath := fmt.Sprintf("/v2/%s/instant/%s/", side, currencyPair)

	var data map[string]string
	data["amount"] = amount.String()
	if clOrdId != "" {
		data["client_order_id"] = clOrdId
	}
	if marginMode != nil {
		data["margin_mode"] = string(*marginMode)
	}
	if leverage != nil {
		data["leverage"] = leverage.String()
	}
	if reduceOnly {
		data["reduce_only"] = "True"
	}

	err = c.authenticatedFormRequest(response, "POST", urlPath, nil, data)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error placing market %s (for %s): %v", side, amount, response.Reason)
		return
	}

	return
}

func (c *HttpClient) V2BuyInstantOrder(currencyPair string, amount decimal.Decimal, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2InstantOrderResponse, err error) {
	return c.v2InstantOrder("buy", currencyPair, amount, clOrdId, marginMode, leverage, reduceOnly)
}

func (c *HttpClient) V2SellInstantOrder(currencyPair string, amount decimal.Decimal, clOrdId string, marginMode *MarginMode, leverage *decimal.Decimal, reduceOnly bool) (response V2InstantOrderResponse, err error) {
	return c.v2InstantOrder("sell", currencyPair, amount, clOrdId, marginMode, leverage, reduceOnly)
}

type MarketSide string

const (
	LONG  MarketSide = "LONG"
	SHORT MarketSide = "SHORT"
)

type V2DerivativesOpenPosition struct {
	Id                        string          `json:"id"`
	Market                    string          `json:"market"`
	MarketType                MarketType      `json:"market_type"`
	MarginMode                MarginMode      `json:"margin_mode"`
	SettlementCurrency        string          `json:"settlement_currency"`
	EntryPrice                decimal.Decimal `json:"entry_price"`
	PnlPercentage             decimal.Decimal `json:"pnl_percentage"`
	PnlRealized               decimal.Decimal `json:"pnl_realized"`
	PnlSettledSinceInception  decimal.Decimal `json:"pnl_settled_since_inception"`
	Leverage                  decimal.Decimal `json:"leverage"`
	Pnl                       decimal.Decimal `json:"pnl"`
	Size                      decimal.Decimal `json:"size"`
	PnlUnrealized             decimal.Decimal `json:"pnl_unrealized"`
	ImpliedLeverage           decimal.Decimal `json:"implied_leverage"`
	InitialMargin             decimal.Decimal `json:"initial_margin"`
	InitialMarginRatio        decimal.Decimal `json:"initial_margin_ratio"`
	CurrentMargin             decimal.Decimal `json:"current_margin"`
	CollateralReserved        decimal.Decimal `json:"collateral_reserved"`
	MaintenanceMargin         decimal.Decimal `json:"maintenance_margin"`
	MaintenanceMarginRatio    decimal.Decimal `json:"maintenance_margin_ratio"`
	EstimatedLiquidationPrice decimal.Decimal `json:"estimated_liquidation_price"`
	EstimatedClosingFeeAmount decimal.Decimal `json:"estimated_closing_fee_amount"`
	MarkPrice                 decimal.Decimal `json:"mark_price"`
	CurrentValue              decimal.Decimal `json:"current_value"`
	EntryValue                decimal.Decimal `json:"entry_value"`
	StrikePrice               decimal.Decimal `json:"strike_price"`
	Side                      MarketSide      `json:"side"`
}

func (c *HttpClient) V2DerivativesOpenPositions(marketSymbol *string) (response []V2DerivativesOpenPosition, err error) {
	urlPath := "/v2/open_positions/"
	if marketSymbol != nil {
		urlPath = fmt.Sprintf("%s%s/", urlPath, *marketSymbol)
	}

	err = c.authenticatedFormRequest(&response, "GET", urlPath, nil, nil)
	if err != nil {
		return
	}

	return response, nil
}

type V2DerivativesOpenPositionRequest struct {
	PositionId string `json:"position_id"`
}

type MarketType string

const (
	Spot      MarketType = "SPOT"
	Perpetual MarketType = "PERPETUAL"
)

type PositionStatus string

const (
	Open              PositionStatus = "OPEN"
	WaitingSettlement PositionStatus = "WAITING_SETTLEMENT"
	Settled           PositionStatus = "SETTLED"
	Liquidating       PositionStatus = "LIQUIDATING"
)

type V2DerivativesOpenPositionResponse struct {
	Id               string                         `json:"id"`
	Market           string                         `json:"market"`
	MarketType       MarketType                     `json:"market_type"`
	MarginMode       MarginMode                     `json:"margin_mode"`
	PnlCurrency      string                         `json:"pnl_currency"`
	EntryPrice       decimal.Decimal                `json:"entry_price"`
	PnlPercentage    decimal.Decimal                `json:"pnl_percentage"`
	PnlRealized      decimal.Decimal                `json:"pnl_realized"`
	PnlSettled       decimal.Decimal                `json:"pnl_settled"`
	Leverage         decimal.Decimal                `json:"leverage"`
	Pnl              decimal.Decimal                `json:"pnl"`
	AmountDelta      decimal.Decimal                `json:"amount_delta"`
	TimeOpened       StringDatetimeFromMicroseconds `json:"time_opened"`
	TimeClosed       StringDatetimeFromMicroseconds `json:"time_closed"`
	Status           PositionStatus                 `json:"status"`
	ExitPrice        decimal.Decimal                `json:"exit_price"`
	SettlementPrice  decimal.Decimal                `json:"settlement_price"`
	ClosingFeeAmount decimal.Decimal                `json:"closing_fee_amount"`
}

func (c *HttpClient) V2DerivativesClosePosition(positionId string) (response V2DerivativesOpenPositionResponse, err error) {
	urlPath := "/v2/close_position/"
	if positionId == "" {
		err = errors.New("positionId is required")
		return
	}
	requestPayload := V2DerivativesOpenPositionRequest{
		PositionId: positionId,
	}

	err = c.authenticatedJsonRequest(&response, "POST", urlPath, nil, requestPayload)
	if err != nil {
		return
	}

	return response, nil
}

type ClosePositionOrderType string

const (
	Market ClosePositionOrderType = "MARKET"
)

type V2DerivativesOpenPositionsRequest struct {
	MarginMode *MarginMode            `json:"margin_mode,omitempty"`
	Market     *string                `json:"market,omitempty"`
	OrderType  ClosePositionOrderType `json:"order_type"`
}
type V2DerivativesOpenPositionsResponse struct {
	Closed []struct {
		Id               string                         `json:"id"`
		Market           string                         `json:"market"`
		MarketType       MarketType                     `json:"market_type"`
		MarginMode       MarginMode                     `json:"margin_mode"`
		PnlCurrency      decimal.Decimal                `json:"pnl_currency"`
		EntryPrice       decimal.Decimal                `json:"entry_price"`
		PnlPercentage    decimal.Decimal                `json:"pnl_percentage"`
		PnlRealized      decimal.Decimal                `json:"pnl_realized"`
		PnlSettled       decimal.Decimal                `json:"pnl_settled"`
		Leverage         decimal.Decimal                `json:"leverage"`
		Pnl              decimal.Decimal                `json:"pnl"`
		AmountDelta      decimal.Decimal                `json:"amount_delta"`
		TimeOpened       StringDatetimeFromMicroseconds `json:"time_opened"`
		TimeClosed       StringDatetimeFromMicroseconds `json:"time_closed"`
		Status           PositionStatus                 `json:"status"`
		ExitPrice        decimal.Decimal                `json:"exit_price"`
		SettlementPrice  decimal.Decimal                `json:"settlement_price"`
		ClosingFeeAmount decimal.Decimal                `json:"closing_fee_amount"`
	} `json:"closed"`
	Failed []struct {
		Id               string                         `json:"id"`
		Market           string                         `json:"market"`
		MarketType       MarketType                     `json:"market_type"`
		MarginMode       MarginMode                     `json:"margin_mode"`
		PnlCurrency      decimal.Decimal                `json:"pnl_currency"`
		EntryPrice       decimal.Decimal                `json:"entry_price"`
		PnlPercentage    decimal.Decimal                `json:"pnl_percentage"`
		PnlRealized      decimal.Decimal                `json:"pnl_realized"`
		PnlSettled       decimal.Decimal                `json:"pnl_settled"`
		Leverage         decimal.Decimal                `json:"leverage"`
		Pnl              decimal.Decimal                `json:"pnl"`
		AmountDelta      decimal.Decimal                `json:"amount_delta"`
		TimeOpened       StringDatetimeFromMicroseconds `json:"time_opened"`
		TimeClosed       StringDatetimeFromMicroseconds `json:"time_closed"`
		Status           PositionStatus                 `json:"status"`
		ExitPrice        decimal.Decimal                `json:"exit_price"`
		SettlementPrice  decimal.Decimal                `json:"settlement_price"`
		ClosingFeeAmount decimal.Decimal                `json:"closing_fee_amount"`
	} `json:"failed"`
}

func (c *HttpClient) V2DerivativesClosePositions(orderType ClosePositionOrderType, marginMode *MarginMode, market *string) (response V2DerivativesOpenPositionsResponse, err error) {
	urlPath := "/v2/close_positions/"
	requestPayload := V2DerivativesOpenPositionsRequest{
		OrderType:  orderType,
		MarginMode: marginMode,
		Market:     market,
	}

	err = c.authenticatedJsonRequest(&response, "POST", urlPath, nil, requestPayload)
	if err != nil {
		return
	}

	return response, nil
}

type V2DerivativesMarginInfoResponse struct {
	AccountMargin          decimal.Decimal `json:"account_margin"`
	AccountMarginAvailable decimal.Decimal `json:"account_margin_available"`
	AccountMarginReserved  decimal.Decimal `json:"account_margin_reserved"`
	Assets                 []struct {
		Asset           string          `json:"asset"`
		Available       decimal.Decimal `json:"available"`
		MarginAvailable decimal.Decimal `json:"margin_available"`
		Reserved        decimal.Decimal `json:"reserved"`
		TotalAmount     decimal.Decimal `json:"total_amount"`
	} `json:"assets"`
	ImpliedLeverage        decimal.Decimal `json:"implied_leverage"`
	InitialMarginRatio     decimal.Decimal `json:"initial_margin_ratio"`
	MaintenanceMarginRatio decimal.Decimal `json:"maintenance_margin_ratio"`
}

func (c *HttpClient) V2DerivativesMarginInfo() (response V2DerivativesMarginInfoResponse, err error) {
	urlPath := "/v2/margin_info/"

	err = c.authenticatedJsonRequest(&response, "GET", urlPath, nil, nil)
	if err != nil {
		return
	}

	return response, nil
}

type Sort string

const (
	Descending Sort = "desc"
	Ascending  Sort = "asc"
)

type V2DerivativesPositionsHistoryListResponse struct {
	Id              string                         `json:"id"`
	Market          string                         `json:"market"`
	MarketType      MarketType                     `json:"market_type"`
	MarginMode      MarginMode                     `json:"margin_mode"`
	PnlCurrency     string                         `json:"pnl_currency"`
	EntryPrice      decimal.Decimal                `json:"entry_price"`
	PnlPercentage   decimal.Decimal                `json:"pnl_percentage"`
	PnlRealized     decimal.Decimal                `json:"pnl_realized"`
	PnlSettled      decimal.Decimal                `json:"pnl_settled"`
	Leverage        decimal.Decimal                `json:"leverage"`
	Pnl             decimal.Decimal                `json:"pnl"`
	AmountDelta     decimal.Decimal                `json:"amount_delta"`
	TimeOpened      StringDatetimeFromMicroseconds `json:"time_opened"`
	TimeClosed      StringDatetimeFromMicroseconds `json:"time_closed"`
	Status          PositionStatus                 `json:"status"`
	ExitPrice       decimal.Decimal                `json:"exit_price"`
	SettlementPrice decimal.Decimal                `json:"settlement_price"`
}

func (c *HttpClient) V2DerivativesPositionsHistoryList(marketSymbol *string, sort *Sort, page *int64, perPage *int64) (response []V2DerivativesPositionsHistoryListResponse, err error) {
	if sort == nil {
		sortValue := Descending
		sort = &sortValue
	}
	if page == nil {
		pageValue := int64(1)
		page = &pageValue
	}
	urlPath := "/v2/position_history/"
	if marketSymbol != nil {
		urlPath += "/" + *marketSymbol + "/"
	}
	urlParams := make(url.Values)
	urlParams.Set("sort", string(*sort))
	urlParams.Set("page", strconv.FormatInt(*page, 10))
	if perPage != nil {
		urlParams.Set("per_page", strconv.FormatInt(*perPage, 10))
	}

	err = c.authenticatedJsonRequest(&response, "GET", urlPath, &urlParams, nil)
	if err != nil {
		return
	}

	return response, nil
}

type SettlementType string

const (
	Periodic SettlementType = "PERIODIC"
	Closed   SettlementType = "CLOSED"
)

type V2DerivativesPositionsSettlementTransactionListResponse struct {
	TransactionId              string                         `json:"transaction_id"`
	PositionId                 string                         `json:"position_id"`
	SettlementTime             StringDatetimeFromMicroseconds `json:"settlement_time"`
	SettlementType             SettlementType                 `json:"settlement_type"`
	SettlementPrice            decimal.Decimal                `json:"settlement_price"`
	Market                     string                         `json:"market"`
	MarketType                 MarketType                     `json:"market_type"`
	PnlCurrency                string                         `json:"pnl_currency"`
	PnlSettled                 decimal.Decimal                `json:"pnl_settled"`
	PnlComponentPrice          decimal.Decimal                `json:"pnl_component_price"`
	PnlComponentFees           decimal.Decimal                `json:"pnl_component_fees"`
	PnlComponentFunding        decimal.Decimal                `json:"pnl_component_funding"`
	PnlComponentSocializedLoss decimal.Decimal                `json:"pnl_component_socialized_loss"`
	MarginMode                 MarginMode                     `json:"margin_mode"`
	Size                       decimal.Decimal                `json:"size"`
	StrikePrice                decimal.Decimal                `json:"strike_price"`
}

func (c *HttpClient) V2DerivativesPositionsSettlementTransactionList(marketTransactionId *string, offset *int64, limit *int64, sort *Sort, sinceTimestamp *int64, untilTimestamp *int64, sinceId *int64) (response []V2DerivativesPositionsSettlementTransactionListResponse, err error) {
	if offset == nil {
		offsetValue := int64(0)
		offset = &offsetValue
	}
	if *offset > 200000 {
		err = errors.New("invalid offset")
		return
	}
	if limit == nil {
		limitValue := int64(100)
		limit = &limitValue
	}
	if *limit > 1000 {
		err = errors.New("invalid limit")
		return
	}
	if sort == nil {
		sortValue := Descending
		sort = &sortValue
	}
	urlPath := "/v2/position_settlement_transactions/"
	if marketTransactionId != nil {
		urlPath += "/" + *marketTransactionId + "/"
	}
	urlParams := make(url.Values)
	urlParams.Set("offset", strconv.FormatInt(*offset, 10))
	urlParams.Set("limit", strconv.FormatInt(*limit, 10))
	urlParams.Set("sort", string(*sort))
	if sinceTimestamp != nil {
		urlParams.Set("since_timestamp", strconv.FormatInt(*sinceTimestamp, 10))
	}
	if untilTimestamp != nil {
		urlParams.Set("until_timestamp", strconv.FormatInt(*untilTimestamp, 10))
	}
	if sinceId != nil {
		urlParams.Set("since_id", strconv.FormatInt(*sinceId, 10))
	}

	err = c.authenticatedJsonRequest(&response, "GET", urlPath, &urlParams, nil)
	if err != nil {
		return
	}

	return response, nil
}

type V2DerivativesAdjustCollateralValueForPositionRequest struct {
	PositionId string          `json:"position_id"`
	NewAmount  decimal.Decimal `json:"new_amount"`
}

type V2DerivativesAdjustCollateralValueForPositionResponse struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (c *HttpClient) V2DerivativesAdjustCollateralValueForPosition(positionId string, newAmount decimal.Decimal) (response V2DerivativesAdjustCollateralValueForPositionResponse, err error) {
	if positionId == "" {
		err = errors.New("positionId is empty")
		return
	}
	requestPayload := V2DerivativesAdjustCollateralValueForPositionRequest{
		PositionId: positionId,
		NewAmount:  newAmount,
	}

	urlPath := "/v2/adjust_position_collateral/"

	err = c.authenticatedJsonRequest(&response, "POST", urlPath, nil, requestPayload)
	if err != nil {
		return
	}

	return response, nil
}

type V2DerivativesCollateralCurrenciesResponse struct {
	Currency string          `json:"currency"`
	Haircut  decimal.Decimal `json:"haircut"`
}

func (c *HttpClient) V2DerivativesCollateralCurrencies() (response []V2DerivativesCollateralCurrenciesResponse, err error) {
	urlPath := "/v2/collateral_currencies/"

	err = c.authenticatedJsonRequest(&response, "GET", urlPath, nil, nil)
	if err != nil {
		return
	}

	return response, nil
}

type V2DerivativesLeverageSettingsListResponse struct {
	LeverageCurrent decimal.Decimal `json:"leverage_current"`
	LeverageMax     decimal.Decimal `json:"leverage_max"`
	MarginMode      MarginMode      `json:"margin_mode"`
	Market          string          `json:"market"`
}

func (c *HttpClient) V2DerivativesLeverageSettingsList(marginMode MarginMode, market string) (response []V2DerivativesLeverageSettingsListResponse, err error) {
	urlPath := "/v2/leverage_settings/"
	urlParams := make(url.Values)
	urlParams.Set("margin_mode", string(marginMode))
	urlParams.Set("market", market)

	err = c.authenticatedJsonRequest(&response, "GET", urlPath, &urlParams, nil)
	if err != nil {
		return
	}

	return response, nil
}

type V2DerivativesUpdateLeverageSettingWithOverrideRequest struct {
	Leverage   decimal.Decimal `json:"leverage"`
	MarginMode MarginMode      `json:"margin_mode"`
	Market     string          `json:"market"`
}

type V2DerivativesUpdateLeverageSettingWithOverrideResponse struct {
	LeverageCurrent decimal.Decimal `json:"leverage_current"`
	LeverageMax     decimal.Decimal `json:"leverage_max"`
	MarginMode      MarginMode      `json:"margin_mode"`
	Market          string          `json:"market"`
}

func (c *HttpClient) V2DerivativesUpdateLeverageSettingWithOverride(leverage decimal.Decimal, marginMode MarginMode, market string) (response V2DerivativesUpdateLeverageSettingWithOverrideResponse, err error) {
	urlPath := "/v2/leverage_settings/"
	requestPayload := V2DerivativesUpdateLeverageSettingWithOverrideRequest{
		Leverage:   leverage,
		MarginMode: marginMode,
		Market:     market,
	}

	err = c.authenticatedJsonRequest(&response, "POST", urlPath, nil, requestPayload)
	if err != nil {
		return
	}

	return response, nil
}

type V2WebsocketsTokenResponse struct {
	Token    string `json:"token"`
	ValidSec uint32 `json:"valid_sec"`
	UserId   uint32 `json:"user_id"`
}

// V2WebsocketsToken generates an ephemeral token, which allows user to subscribe to private
// websocket events. These events include ClientOrderIds (and potentially additional private data)
func (c *HttpClient) V2WebsocketsToken() (response V2WebsocketsTokenResponse, err error) {
	err = c.authenticatedFormRequest(&response, "POST", "/v2/websockets_token/", nil, nil)
	if err != nil {
		return
	}
	return
}
