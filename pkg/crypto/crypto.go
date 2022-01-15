package crypto

type Crypto interface {
	FindSymbol(symbol string) (Symbol, bool)
	GetSymbols() []Symbol
	GetCurrencies() []string
	GetQuotes(currency string, symbol ...string) ([]Quote, error)
	GetOHLCV(currency string, symbol ...string) ([]Ohlcv, error)
}
