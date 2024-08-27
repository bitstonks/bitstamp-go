package main

import (
	"fmt"
	"log"

	"github.com/bitstonks/bitstamp-go/pkg/http"
	"github.com/shopspring/decimal"
)

func main() {
	api := http.NewHttpClient(
		http.Credentials("invalid", "invalid"),
	)
	currencyPair := "btcusd-perp"
	market := "BTC/USD-PERP"

	// public endpoints
	ticker1, err := api.V2Ticker(currencyPair)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("TICKER: %+v\n", ticker1)

	ticker2, err := api.V2HourlyTicker(currencyPair)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("HOURLY TICKER: %+v\n", ticker2)

	ob, err := api.V2OrderBook(currencyPair, 2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("ORDER BOOK - HIGHEST BID: %+v\n", ob.Bids[0])
	fmt.Printf("ORDER BOOK - LOWEST ASK: %+v\n", ob.Asks[0])

	txs, err := api.V2Transactions(currencyPair, "hour")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("TRANSACTIONS: %+v\n", txs)

	info, err := api.V2TradingPairsInfo()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("TRADING PAIRS: %+v\n", info[0])

	ohlc, err := api.V2Ohlc(currencyPair, 60, 2, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("CANDLES: %+v\n", ohlc.Data.Candles)

	eurusd, err := api.V2EurUsd()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("EURUSD: %+v\n", eurusd)

	balances, err := api.V2AccountBalances()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("BALANCES: %+v\n", balances)

	openOrders, err := api.V2OpenOrders(currencyPair)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("OPEN ORDERS: %+v\n", openOrders)

	openPositions, err := api.V2DerivativesOpenPositions(&currencyPair)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("OPEN POSITIONS: %+v\n", openPositions)

	if len(openPositions) == 0 {
		amount := decimal.NewFromFloat(0.1)
		marginMode := http.Isolated
		leverage := decimal.NewFromInt(3)
		order, err := api.V2BuyMarketOrder(currencyPair, amount, "", &marginMode, &leverage, false)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("MARKET ORDER: %+v\n", order)
		openPositions, err := api.V2DerivativesOpenPositions(&currencyPair)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("OPEN POSITIONS: %+v\n", openPositions)
		for _, position := range openPositions {
			adjustPosition, err := api.V2DerivativesAdjustCollateralValueForPosition(position.Id, position.CollateralReserved.Add(decimal.NewFromInt(1)))
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("ADJUST COLLATERAL VALUE: %+v\n", adjustPosition)

			closePosition, err := api.V2DerivativesClosePosition(position.Id)
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("CLOSE POSITION: %+v\n", closePosition)
		}
	}

	if len(openPositions) > 0 {
		closePositions, err := api.V2DerivativesClosePositions(http.Market, nil, &market)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("CLOSE POSITIONS: %+v\n", closePositions)
	}

	positionHistory, err := api.V2DerivativesPositionsHistoryList(&currencyPair, nil, nil, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("POSITION HISTORY: %+v\n", positionHistory)

	collateralCurrencies, err := api.V2DerivativesCollateralCurrencies()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("COLLETERAL CURRENCIES: %+v\n", collateralCurrencies)

	leverageSettingsList, err := api.V2DerivativesLeverageSettingsList(http.Cross, market)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("LEVERAGE SETTINGS: %+v\n", leverageSettingsList)

	marginInfo, err := api.V2DerivativesMarginInfo()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("MARGIN INFO: %+v\n", marginInfo)

	positionSettlementTransactionList, err := api.V2DerivativesPositionsSettlementTransactionList(nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("POSITION SETTLEMENT TRANSACTION LIST: %+v\n", positionSettlementTransactionList)
}
