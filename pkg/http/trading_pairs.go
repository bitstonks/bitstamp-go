package http

type rounding struct {
	Base    int32
	Counter int32
}

// to generate:
// curl -s https://www.bitstamp.net/api/v2/trading-pairs-info/ | jq -r '.[] | "\"\(.url_symbol)\": {\(.base_decimals), \(.counter_decimals)},"' | sort
var roundings = map[string]rounding{
	"1incheur":  {2, 2},
	"1inchusd":  {2, 2},
	"aavebtc":   {8, 8},
	"aaveeur":   {8, 2},
	"aaveusd":   {8, 2},
	"adabtc":    {8, 8},
	"adaeur":    {8, 5},
	"adausd":    {8, 5},
	"algobtc":   {8, 8},
	"algoeur":   {8, 5},
	"algousd":   {8, 5},
	"alphaeur":  {8, 5},
	"alphausd":  {8, 5},
	"ampeur":    {8, 5},
	"ampusd":    {8, 5},
	"anteur":    {1, 2},
	"antusd":    {1, 2},
	"apeeur":    {2, 2},
	"apeusd":    {2, 2},
	"audiobtc":  {8, 8},
	"audioeur":  {8, 5},
	"audiousd":  {8, 5},
	"avaxeur":   {8, 5},
	"avaxusd":   {8, 5},
	"axseur":    {8, 5},
	"axsusd":    {8, 5},
	"bandeur":   {2, 3},
	"bandusd":   {2, 3},
	"bateur":    {8, 5},
	"batusd":    {8, 5},
	"bchbtc":    {8, 8},
	"bcheur":    {8, 2},
	"bchusd":    {8, 2},
	"blureur":   {3, 4},
	"blurusd":   {3, 4},
	"btceur":    {8, 0},
	"btcgbp":    {8, 0},
	"btcusd":    {8, 0},
	"btcusdc":   {8, 0},
	"btcusdt":   {8, 0},
	"chzeur":    {8, 5},
	"chzusd":    {8, 5},
	"compeur":   {8, 2},
	"compusd":   {8, 2},
	"crveur":    {8, 5},
	"crvusd":    {8, 5},
	"cspreur":   {1, 5},
	"csprusd":   {1, 5},
	"ctsieur":   {1, 4},
	"ctsiusd":   {1, 4},
	"cvxeur":    {3, 2},
	"cvxusd":    {3, 2},
	"daiusd":    {5, 5},
	"dgldeur":   {5, 1},
	"dgldusd":   {5, 1},
	"dogeeur":   {2, 5},
	"dogeusd":   {2, 5},
	"doteur":    {2, 3},
	"dotusd":    {2, 3},
	"dydxeur":   {8, 3},
	"dydxusd":   {8, 3},
	"enjeur":    {8, 5},
	"enjusd":    {8, 5},
	"enseur":    {2, 2},
	"ensusd":    {2, 2},
	"eth2eth":   {8, 8},
	"ethbtc":    {8, 8},
	"etheur":    {8, 1},
	"ethgbp":    {8, 1},
	"ethusd":    {8, 1},
	"ethusdc":   {8, 1},
	"ethusdt":   {8, 1},
	"eurcveur":  {2, 4},
	"eurcvusdt": {2, 4},
	"euroceur":  {2, 4},
	"eurocusdc": {2, 4},
	"eurteur":   {5, 5},
	"eurtusd":   {5, 5},
	"eurusd":    {5, 5},
	"feteur":    {8, 5},
	"fetusd":    {8, 5},
	"flreur":    {1, 5},
	"flrusd":    {1, 5},
	"ftmeur":    {8, 5},
	"ftmusd":    {8, 5},
	"galaeur":   {8, 5},
	"galausd":   {8, 5},
	"gbpusd":    {5, 5},
	"godseur":   {2, 2},
	"godsusd":   {2, 2},
	"grteur":    {8, 5},
	"grtusd":    {8, 5},
	"gusdusd":   {5, 5},
	"hbareur":   {8, 5},
	"hbarusd":   {8, 5},
	"imxeur":    {2, 2},
	"imxusd":    {2, 2},
	"injeur":    {2, 3},
	"injusd":    {2, 3},
	"knceur":    {8, 5},
	"kncusd":    {8, 5},
	"ldoeur":    {2, 4},
	"ldousd":    {2, 4},
	"linkbtc":   {8, 8},
	"linkeur":   {8, 2},
	"linkgbp":   {8, 2},
	"linkusd":   {8, 2},
	"lmwreur":   {3, 4},
	"lmwrusd":   {3, 4},
	"lrceur":    {1, 4},
	"lrcusd":    {1, 4},
	"ltcbtc":    {8, 8},
	"ltceur":    {8, 2},
	"ltcgbp":    {8, 2},
	"ltcusd":    {8, 2},
	"manaeur":   {2, 2},
	"manausd":   {2, 2},
	"maticeur":  {8, 5},
	"maticusd":  {8, 5},
	"mkreur":    {8, 2},
	"mkrusd":    {8, 2},
	"mpleur":    {3, 2},
	"mplusd":    {3, 2},
	"neareur":   {2, 3},
	"nearusd":   {2, 3},
	"nexoeur":   {2, 2},
	"nexousd":   {2, 2},
	"paxusd":    {5, 5},
	"pepeeur":   {1, 8},
	"pepeusd":   {1, 8},
	"perpeur":   {8, 3},
	"perpusd":   {8, 3},
	"pyusdeur":  {2, 4},
	"pyusdusd":  {2, 4},
	"radeur":    {2, 2},
	"radusd":    {2, 2},
	"rlyeur":    {0, 4},
	"rlyusd":    {0, 4},
	"rndreur":   {2, 3},
	"rndrusd":   {2, 3},
	"sandeur":   {8, 5},
	"sandusd":   {8, 5},
	"sgbeur":    {8, 5},
	"sgbusd":    {8, 5},
	"shibeur":   {0, 8},
	"shibusd":   {0, 8},
	"skleur":    {8, 5},
	"sklusd":    {8, 5},
	"slpeur":    {0, 5},
	"slpusd":    {0, 5},
	"snxeur":    {8, 5},
	"snxusd":    {8, 5},
	"soleur":    {2, 4},
	"solusd":    {2, 4},
	"storjeur":  {8, 5},
	"storjusd":  {8, 5},
	"suieur":    {4, 3},
	"suiusd":    {4, 3},
	"sushieur":  {8, 5},
	"sushiusd":  {8, 5},
	"traceur":   {2, 4},
	"tracusd":   {2, 4},
	"umaeur":    {8, 2},
	"umausd":    {8, 2},
	"unibtc":    {8, 8},
	"unieur":    {8, 5},
	"uniusd":    {8, 5},
	"usdceur":   {5, 5},
	"usdcusd":   {5, 5},
	"usdcusdt":  {5, 5},
	"usdteur":   {5, 5},
	"usdtusd":   {5, 5},
	"vchfeur":   {2, 4},
	"vchfusd":   {2, 4},
	"vegaeur":   {2, 3},
	"vegausd":   {2, 3},
	"veureur":   {2, 4},
	"veurusd":   {2, 4},
	"vexteur":   {3, 4},
	"vextusd":   {3, 4},
	"wbtcbtc":   {4, 4},
	"wecaneur":  {2, 5},
	"wecanusd":  {2, 5},
	"xlmbtc":    {8, 8},
	"xlmeur":    {8, 5},
	"xlmgbp":    {8, 5},
	"xlmusd":    {8, 5},
	"xrpbtc":    {8, 8},
	"xrpeur":    {8, 5},
	"xrpgbp":    {8, 5},
	"xrpusd":    {8, 5},
	"xrpusdt":   {8, 5},
	"yfieur":    {8, 2},
	"yfiusd":    {8, 2},
	"zrxeur":    {8, 5},
	"zrxusd":    {8, 5},
}
