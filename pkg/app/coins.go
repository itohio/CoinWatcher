package app

import (
	"encoding/json"

	"github.com/itohio/CoinWatcher/pkg/crypto"
	"github.com/itohio/CoinWatcher/pkg/widgets/coin"
)

type Coin struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Coins struct {
	Coins []Coin `json:"coins"`
}

func (a *App) defaultCoins() {
	if a.feed == nil {
		return
	}
	for _, ss := range a.feed.GetSymbols() {
		switch ss.Symbol {
		case "BTC":
			fallthrough
		case "ETH":
			fallthrough
		case "DOT":
			fallthrough
		case "ADA":
			a.addSymbol(ss)
		}
	}
}

func (a *App) loadCoins() {
	reader, err := a.reader("coins.json")
	if err != nil {
		a.defaultCoins()
		return
	}
	defer reader.Close()

	var coins Coins

	err = json.NewDecoder(reader).Decode(&coins)
	if err != nil {
		a.defaultCoins()
		return
	}

	a.Lock()
	a.coinData = a.coinData[:0]
	a.data.Reload()
	a.Unlock()

	for _, c := range coins.Coins {
		a.addSymbol(crypto.Symbol{
			Name:   c.Name,
			Symbol: c.Symbol,
		})
	}
}

func (a *App) saveCoins() {
	writer, err := a.writer("coins.json")
	if err != nil {
		return
	}
	defer writer.Close()

	a.Lock()
	defer a.Unlock()

	coins := Coins{
		Coins: make([]Coin, len(a.coinData)),
	}

	for i, c := range a.coinData {
		if cn, ok := c.(*coin.CoinData); ok {
			coins.Coins[i] = Coin{
				Symbol: cn.Symbol.Symbol,
				Name:   cn.Symbol.Name,
			}
		}
	}

	data, err := json.Marshal(&coins)
	if err != nil {
		return
	}

	if _, err := writer.Write(data); err != nil {

	}
}
