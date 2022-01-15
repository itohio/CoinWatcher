package crypto

import (
	"fmt"
	"sort"
	"strings"

	cmc "github.com/hexoul/go-coinmarketcap"
	"github.com/hexoul/go-coinmarketcap/types"
	"github.com/itohio/CoinWatcher/pkg/logger"
)

type coinmarketcap struct {
	client     *cmc.Client
	key        string
	symbols    []Symbol
	currencies []string
	iconCache  Cache
}

var _ Crypto = &coinmarketcap{}

func NewCMC(key string, iconCache Cache) *coinmarketcap {
	client := cmc.GetInstanceWithKey(key)

	list, err := client.CryptoListingsLatest(&types.Options{
		Limit: 1500,
	})
	if err != nil {
		panic(fmt.Errorf("Could not get symbol list: %v", err))
	}
	symbols := make([]Symbol, len(list.CryptoMarket))
	sStr := make([]string, len(list.CryptoMarket))
	currencies := make(map[string]struct{})
	for i, l := range list.CryptoMarket {
		symbols[i] = Symbol{
			Id:        l.ID,
			Name:      l.Name,
			Symbol:    l.Symbol,
			iconCache: iconCache,
		}
		sStr[i] = fmt.Sprint(l.ID)

		for q := range l.Quote {
			currencies[q] = struct{}{}
		}
	}

	currencyList := make([]string, 0, len(currencies))
	for q := range currencies {
		currencyList = append(currencyList, q)
	}
	sort.Sort(sort.StringSlice(currencyList))

	logger.Log.Debug().Int("symbols", len(sStr)).Int("currencies", len(currencies)).Msg("Fetched symbols")

	info, err := cmcLoadSymbols(sStr...)
	if err != nil {
		panic(fmt.Errorf("Could not get symbols info: %v", err))
	}

	for i, s := range sStr {
		if sInfo, ok := info.CryptoInfo[s]; ok {
			symbols[i].IconURL = sInfo.Logo
		}
	}

	return &coinmarketcap{
		client:     client,
		key:        key,
		symbols:    symbols,
		currencies: currencyList,
		iconCache:  iconCache,
	}
}

func cmcLoadSymbols(symbols ...string) (info *types.CryptoInfoMap, err error) {
	const N = 1000
	client := cmc.GetInstance()
	var slice []string
	for len(symbols) > 0 {
		n := N
		if len(symbols) < n {
			n = len(symbols)
		}
		slice, symbols = symbols[:n], symbols[n:]
		var infoTmp *types.CryptoInfoMap
		infoTmp, err = client.CryptoInfo(&types.Options{
			ID: strings.Join(slice, ","),
		})
		if err != nil {
			logger.Log.Error().Err(err).Str("slice", strings.Join(slice, ",")).Msg("Failed CryptoInfo")
			return
		}

		if info == nil {
			info = infoTmp
			continue
		}

		for k, v := range infoTmp.CryptoInfo {
			info.CryptoInfo[k] = v
		}
	}
	logger.Log.Debug().Msg("Loaded symbols")

	return
}

func (c *coinmarketcap) GetSymbols() []Symbol {
	return c.symbols
}

func (c *coinmarketcap) GetCurrencies() []string {
	return c.currencies
}

func (c *coinmarketcap) FindSymbol(symbol string) (Symbol, bool) {
	for _, s := range c.symbols {
		if s.Symbol == symbol {
			return s, true
		}
	}
	return Symbol{}, false
}

func cmcQ2Q(s Symbol, q types.Quote) Quote {
	ret := Quote{
		Symbol: s,
	}

	ret.MarketCap = q.MarketCap
	ret.PercentChange1H = q.PercentChange1H
	ret.PercentChange24H = q.PercentChange24H
	ret.PercentChange30D = q.PercentChange30D
	ret.PercentChange7D = q.PercentChange7D

	ret.Price = q.Price
	ret.Volume24Hbase = q.Volume24Hbase
	ret.Volume24Hquote = q.Volume24Hquote
	ret.Volume24H = q.Volume24H
	ret.Volume7D = q.Volume7D
	ret.Volume30D = q.Volume30D

	//ret.LastUpdated = time.Parse()

	return ret
}

func (c *coinmarketcap) GetQuotes(currency string, symbol ...string) ([]Quote, error) {
	qts, err := c.client.CryptoMarketQuotesLatest(&types.Options{
		Symbol: strings.Join(symbol, ","),
	})
	if err != nil {
		logger.Log.Error().Err(err).Str("symbol", strings.Join(symbol, ",")).Msg("Could not get latest quotes")
		return nil, err
	}

	quotes := make([]Quote, 0, len(qts.CryptoMarket))
	for s, q := range qts.CryptoMarket {
		if quote, ok := q.Quote[currency]; ok {
			if sym, ok := c.FindSymbol(s); ok {
				quotes = append(quotes, cmcQ2Q(sym, *quote))
			}
		}
	}

	return quotes, nil
}

func (c *coinmarketcap) GetOHLCV(currency string, symbol ...string) ([]Ohlcv, error) {
	return nil, fmt.Errorf("not implemented")
}
