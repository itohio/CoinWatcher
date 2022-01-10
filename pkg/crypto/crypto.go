package crypto

import (
	"fmt"
	"image"
	"net/http"
	"sort"
	"strings"

	cmc "github.com/hexoul/go-coinmarketcap"
	"github.com/hexoul/go-coinmarketcap/types"
	"github.com/itohio/CoinWatcher/pkg/logger"
)

type Symbol struct {
	iconCache Cache
	Name      string
	Symbol    string
	IconURL   string
}

type Quote struct {
	Symbol
	types.Quote
}

type Crypto struct {
	client     *cmc.Client
	key        string
	symbols    []Symbol
	currencies []string
	iconCache  Cache
}

func New(key string, iconCache Cache) *Crypto {
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
			Name:      l.Name,
			Symbol:    l.Symbol,
			iconCache: iconCache,
		}
		sStr[i] = l.Symbol

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

	info, err := LoadSymbols(sStr...)
	if err != nil {
		panic(fmt.Errorf("Could not get symbols info: %v", err))
	}

	for i, s := range sStr {
		if sInfo, ok := info.CryptoInfo[s]; ok {
			symbols[i].IconURL = sInfo.Logo
		}
	}

	return &Crypto{
		client:     client,
		key:        key,
		symbols:    symbols,
		currencies: currencyList,
		iconCache:  iconCache,
	}
}

func LoadSymbols(symbols ...string) (info *types.CryptoInfoMap, err error) {
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
			Symbol: strings.Join(slice, ","),
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

	return
}

func (c *Crypto) GetSymbols() []Symbol {
	return c.symbols
}

func (c *Crypto) GetCurrencies() []string {
	return c.currencies
}

func (c *Crypto) findSymbol(symbol string) (Symbol, bool) {
	for _, s := range c.symbols {
		if s.Symbol == symbol {
			return s, true
		}
	}
	return Symbol{}, false
}

func (c *Crypto) GetQuotes(currency string, symbol ...string) ([]Quote, error) {
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
			if sym, ok := c.findSymbol(s); ok {
				quotes = append(quotes, Quote{
					Symbol: sym,
					Quote:  *quote,
				})
			}
		}
	}

	return quotes, nil
}

func (s *Symbol) Icon() image.Image {
	if s.iconCache != nil {
		if img, err := s.iconCache.LoadImage(s.IconURL); err == nil {
			return img
		}
	}

	response, err := http.Get(s.IconURL)
	if err != nil {
		logger.Log.Warn().Str("url", s.IconURL).Err(err).Msg("No icon")
		return image.NewGray(image.Rectangle{Max: image.Point{32, 32}})
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		logger.Log.Warn().Str("url", s.IconURL).Err(err).Msg("No icon")
		return image.NewGray(image.Rectangle{Max: image.Point{32, 32}})
	}

	logger.Log.Debug().Str("url", s.IconURL).Float32("w", float32(img.Bounds().Dx())).Float32("h", float32(img.Bounds().Dy())).Msg("icon img")

	if s.iconCache != nil {
		s.iconCache.SaveImage(s.IconURL, img)
	}

	return img
}
