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

type rounding struct {
	Base    int32
	Counter int32
}

// to generate:
// curl -s https://www.bitstamp.net/api/v2/trading-pairs-info/ | jq -r '.[] | "\"\(.url_symbol)\": {\(.base_decimals), \(.counter_decimals)},"' | sort
var roundings = map[string]rounding{
	"aavebtc":  {8, 8},
	"aaveeur":  {8, 2},
	"aaveusd":  {8, 2},
	"algobtc":  {8, 8},
	"algoeur":  {8, 5},
	"algousd":  {8, 5},
	"audiobtc": {8, 8},
	"audioeur": {8, 5},
	"audiousd": {8, 5},
	"batbtc":   {8, 8},
	"bateur":   {8, 5},
	"batusd":   {8, 5},
	"bchbtc":   {8, 8},
	"bcheur":   {8, 2},
	"bchgbp":   {8, 2},
	"bchusd":   {8, 2},
	"btceur":   {8, 2},
	"btcgbp":   {8, 2},
	"btcpax":   {8, 2},
	"btcusd":   {8, 2},
	"btcusdc":  {8, 2},
	"btcusdt":  {8, 2},
	"compbtc":  {8, 8},
	"compeur":  {8, 2},
	"compusd":  {8, 2},
	"crvbtc":   {8, 8},
	"crveur":   {8, 5},
	"crvusd":   {8, 5},
	"daiusd":   {5, 5},
	"eth2eth":  {8, 8},
	"ethbtc":   {8, 8},
	"etheur":   {8, 2},
	"ethgbp":   {8, 2},
	"ethpax":   {8, 2},
	"ethusd":   {8, 2},
	"ethusdc":  {8, 2},
	"ethusdt":  {8, 2},
	"eurusd":   {5, 5},
	"gbpeur":   {5, 5},
	"gbpusd":   {5, 5},
	"grteur":   {8, 5},
	"grtusd":   {8, 5},
	"gusdusd":  {5, 5},
	"kncbtc":   {8, 8},
	"knceur":   {8, 5},
	"kncusd":   {8, 5},
	"linkbtc":  {8, 8},
	"linketh":  {8, 8},
	"linkeur":  {8, 2},
	"linkgbp":  {8, 2},
	"linkusd":  {8, 2},
	"ltcbtc":   {8, 8},
	"ltceur":   {8, 2},
	"ltcgbp":   {8, 2},
	"ltcusd":   {8, 2},
	"mkrbtc":   {8, 8},
	"mkreur":   {8, 2},
	"mkrusd":   {8, 2},
	"omgbtc":   {8, 8},
	"omgeur":   {8, 2},
	"omggbp":   {8, 2},
	"omgusd":   {8, 2},
	"paxeur":   {5, 5},
	"paxgbp":   {5, 5},
	"paxusd":   {5, 5},
	"snxbtc":   {8, 8},
	"snxeur":   {8, 5},
	"snxusd":   {8, 5},
	"umabtc":   {8, 8},
	"umaeur":   {8, 2},
	"umausd":   {8, 2},
	"unibtc":   {8, 8},
	"unieur":   {8, 5},
	"uniusd":   {8, 5},
	"usdceur":  {5, 5},
	"usdcusd":  {5, 5},
	"usdcusdt": {5, 5},
	"usdteur":  {5, 5},
	"usdtusd":  {5, 5},
	"xlmbtc":   {8, 8},
	"xlmeur":   {8, 5},
	"xlmgbp":   {8, 5},
	"xlmusd":   {8, 5},
	"xrpbtc":   {8, 8},
	"xrpeur":   {8, 5},
	"xrpgbp":   {8, 5},
	"xrppax":   {8, 5},
	"xrpusd":   {8, 5},
	"xrpusdt":  {8, 5},
	"yfibtc":   {8, 8},
	"yfieur":   {8, 2},
	"yfiusd":   {8, 2},
	"zrxbtc":   {8, 8},
	"zrxeur":   {8, 5},
	"zrxusd":   {8, 5},
}

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
	msg := xAuth + method + strings.TrimPrefix(strings.TrimPrefix(url_, "https://"), "http://")
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
func (c *ApiClient) V2UserTransactions(currencyPairOrAll string) (response []V2UserTransactionsResponse, err error) {
	if currencyPairOrAll == "all" {
		err = c.authenticatedPostRequest(&response, "/v2/user_transactions/", [2]string{"limit", "1000"})
	} else {
		err = c.authenticatedPostRequest(&response, fmt.Sprintf("/v2/user_transactions/%s/", currencyPairOrAll), [2]string{"limit", "1000"})
	}

	return
}

// POST https://www.bitstamp.net/api/v2/crypto-transactions/
func (c *ApiClient) V2CryptoTransactions() (response []V2CryptoTransactionsResponse, err error) {
	err = c.authenticatedPostRequest(&response, "v2/crypto-transactions/", [2]string{"limit", "1000"})
	if err != nil {
		return
	}
	return
}

type V2CryptoTransactionsResponse struct {
	Currency           string `json:"currency"`
	DestinationAddress string `json:"destination_address"`
	TxID               string `json:"tx_id"`
	Amount             string `json:"amount"`
	DateTime           string `json:"date_time"`
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
func (c *ApiClient) V2OpenOrders(currencyPairOrAll string) (response []V2OpenOrdersResponse, err error) {
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

// POST https://www.bitstamp.net/api/{token}_withdrawal
func (c *ApiClient) V2CryptoWithdrawals(token, address *string, amount *float64,
	memoID, destinationTag string) (response V2CryptoWithdrawalResponse, err error) {
	params := make([][2]string, 0)
	params = append(params, [2]string{"amount", fmt.Sprintf("%d", amount)})
	params = append(params, [2]string{"address", *address})
	if memoID != "" {
		params = append(params, [2]string{"memo_id", memoID})
	}
	if destinationTag != "" {
		params = append(params, [2]string{"destination_tag", memoID})
	}

	err = c.authenticatedPostRequest(&response, fmt.Sprintf("v2/%s_withdrawal/", *token), params...)
	if err != nil {
		return
	}

	return
}

type V2CryptoWithdrawalResponse struct {
	WithdrawalID int64 `json:"withdrawal_id"`
}

// POST https://www.bitstamp.net/api/v2/{token}_address/
func (c *ApiClient) V2TokenDepositAddress(token *string) (response V2TokenDepositAddres, err error) {

	err = c.authenticatedPostRequest(&response, fmt.Sprintf("v2/%s_address/", *token))
	if err != nil {
		return
	}

	return
}

type V2TokenDepositAddres struct {
	Address string `json:"address"`
}

type V2OpenBankWithdrawalRequest struct {
	Amount          float64 `json:"amount"`
	AccountCurrency string  `json:"account_currency"`
	Name            string  `json:"name"`
	Iban            string  `json:"iban"`
	Bic             string  `json:"bic"`
	Address         string  `json:"address"`
	PostalCode      string  `json:"postal_code"`
	City            string  `json:"city"`
	Country         string  `json:"country"`
	Type            string  `json:"type"`
	BankName        string  `json:"bank_name"`
	BankAddress     string  `json:"bank_address"`
	BankPostalCode  string  `json:"bank_postal_code"`
	BankCity        string  `json:"bank_city"`
	BankCountry     string  `json:"bank_country"`
	Currency        string  `json:"currency"`
}

type V2OpenBankWithdrawalResponse struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// POST https://www.bitstamp.net/api/v2/withdrawal/open/
func (c *ApiClient) V2OpenBankWithdrawal(req V2OpenBankWithdrawalRequest) (response V2OpenBankWithdrawalResponse, err error) {
	params := make([][2]string, 0)
	params = append(params, [2]string{"amount", fmt.Sprintf("%d", req.Amount)})
	params = append(params, [2]string{"account_currency", req.AccountCurrency})
	params = append(params, [2]string{"name", req.Name})
	params = append(params, [2]string{"iban", req.Iban})
	params = append(params, [2]string{"bic", req.Bic})
	params = append(params, [2]string{"address", req.Address})
	params = append(params, [2]string{"postal_code", req.PostalCode})
	params = append(params, [2]string{"city", req.City})
	params = append(params, [2]string{"country", req.Country})
	params = append(params, [2]string{"type", req.Type})
	if req.Type == "international" {
		params = append(params, [2]string{"bank_name", req.BankName})
		params = append(params, [2]string{"bank_address", req.BankAddress})
		params = append(params, [2]string{"bank_postal_code", req.BankPostalCode})
		params = append(params, [2]string{"bank_city", req.BankCity})
		params = append(params, [2]string{"bank_country", req.BankCountry})
		params = append(params, [2]string{"currency", req.Currency})
	}

	err = c.authenticatedPostRequest(&response, "/v2/withdrawal/open/", params...)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error: %v", response.Reason)
	}

	return
}

//POST https://www.bitstamp.net/api/v2/withdrawal/status/
func (c *ApiClient) V2BankWithdrawalStatus(withdrawalID int64) (response V2BankWithdrawalStatusResponse, err error) {
	params := make([][2]string, 0)
	params = append(params, [2]string{"id", fmt.Sprintf("%d", withdrawalID)})


	err = c.authenticatedPostRequest(&response, "/v2/withdrawal/status/", params...)
	if err != nil {
		return
	}

	if response.Status == "error" {
		err = fmt.Errorf("error: %v", response.Reason)
	}

	return
}

type V2BankWithdrawalStatusResponse struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// POST https://www.bitstamp.net/api/v2/order_status/
func (c *ApiClient) V2OrderStatus(orderId int64, clOrdId string, omitTx bool) (response V2OrderStatusResponse, err error) {
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

func (c *ApiClient) V2CancelOrder(orderId int64) (response V2CancelOrderResponse, err error) {
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

func (c *ApiClient) v2LimitOrder(side, currencyPair string, price, amount, limitPrice decimal.Decimal, dailyOrder, iocOrder bool, clOrdId string) (response V2LimitOrderResponse, err error) {
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
