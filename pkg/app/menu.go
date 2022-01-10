package app

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/CoinWatcher/pkg/crypto"
	"github.com/itohio/CoinWatcher/pkg/logger"
	"github.com/itohio/CoinWatcher/pkg/widgets/coin"
)

func (a *App) makeMenu() *fyne.Container {
	menu := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			dialog.ShowConfirm(
				"Save coins",
				"You are about to save the coin list.\nAre you sure?",
				func(b bool) {
					if b {
						a.saveCoins()
					}
				},
				a.window,
			)
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			dialog.ShowConfirm(
				"Reload coins",
				"You are about to reload the coin list. All changes will be lost.\nAre you sure?",
				func(b bool) {
					if b {
						a.loadCoins()
					}
				},
				a.window,
			)
		}),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			a.addNewSymbol()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			a.showSettings()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			a.updateQuotes()
		}),
	)

	a.currencyWidget = widget.NewSelect([]string{}, func(s string) {
		a.currency = s
		a.updateQuotes()
		a.saveSettings()
	})
	a.updateCurrencies()

	return container.NewBorder(nil, nil, nil, a.currencyWidget, menu)
}

func (a *App) addNewSymbol() {
	if a.feed == nil {
		dialog.ShowInformation("Add Coin", "Please connect to Coinmarketcap.", a.window)
		return
	}

	symbols := a.feed.GetSymbols()
	options := make([]string, 0, len(symbols))
	for _, s := range symbols {
		if _, ok := a.getSymbol(s.Symbol); ok {
			continue
		}
		options = append(options, fmt.Sprintf("%s (%s)", s.Symbol, s.Name))
	}

	symbolSelect := widget.NewSelectEntry(options)

	dialog.ShowForm(
		"Add a coin",
		"OK",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Coin", symbolSelect),
		},
		func(b bool) {
			if !b {
				return
			}
			if s, ok := a.lookupSymbol(symbolSelect.Text); ok {
				a.addSymbol(s)
			} else {
				dialog.ShowError(fmt.Errorf("Could not find such a coin: %s", symbolSelect.Text), a.window)
			}
		},
		a.window,
	)
}

func (a *App) showSettings() {
	const NOPTS = 10
	options := [NOPTS]string{
		"5 Minutes",
		"10 Minutes",
		"15 Minutes",
		"30 Minutes",
		"1 Hour",
		"2 Hours",
		"4 Hours",
		"8 Hours",
		"12 Hours",
		"24 Hours",
	}
	optionsInt := [NOPTS]time.Duration{
		time.Minute * 5,
		time.Minute * 10,
		time.Minute * 15,
		time.Minute * 30,
		time.Hour,
		time.Hour * 2,
		time.Hour * 4,
		time.Hour * 8,
		time.Hour * 12,
		time.Hour * 24,
	}
	apiKey := widget.NewEntry()
	apiKey.Text = a.apiKey
	interval := widget.NewSelect(options[:], nil)

	for i := range options {
		if a.interval >= optionsInt[NOPTS-i-1] {
			interval.SetSelectedIndex(NOPTS - i - 1)
			break
		}
	}

	dialog.ShowForm(
		"Settings",
		"Save",
		"Discard",
		[]*widget.FormItem{
			widget.NewFormItem("API Key", apiKey),
			widget.NewFormItem("Refresh interval", interval),
		},
		func(b bool) {
			if !b {
				return
			}
			a.apiKey = apiKey.Text
			if interval.SelectedIndex() >= 0 {
				a.interval = optionsInt[interval.SelectedIndex()]
			}
			a.saveSettings()
			a.pbWidget.Refresh()
		},
		a.window,
	)
}

func (a *App) updateCurrencies() {
	if a.feed == nil {
		return
	}
	a.currencyWidget.Options = a.feed.GetCurrencies()
	a.currencyWidget.SetSelected(a.currency)
}

func (a *App) updateQuotes() {
	if a.feed == nil {
		return
	}

	symbols := make([]string, 0, len(a.coinData))
	for _, cd := range a.coinData {
		if s, ok := cd.(*coin.CoinData); ok {
			symbols = append(symbols, s.Symbol.Symbol)
		}
	}
	quotes, err := a.feed.GetQuotes(a.currency, symbols...)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Could not get quotes")
	}

	for _, quote := range quotes {
		a.updateQuote(quote)
	}

	a.lastUpdated = time.Now()
}

func (a *App) delSymbol(symbol string) {
	defer a.data.Reload()
	a.Lock()
	defer a.Unlock()

	var delList []int

	for idx, cd := range a.coinData {
		if coin, ok := cd.(*coin.CoinData); ok {
			if coin.Symbol.Symbol == symbol {
				delList = append(delList, idx)
			}
		}
	}

	for _, idx := range delList {
		a.coinData = append(a.coinData[:idx], a.coinData[idx+1:]...)
	}
}

func (a *App) lookupSymbol(symbol string) (crypto.Symbol, bool) {
	if a.feed == nil {
		return crypto.Symbol{}, false
	}

	for _, s := range a.feed.GetSymbols() {
		if s.Symbol == symbol {
			return s, true
		}
	}
	return crypto.Symbol{}, false
}

func (a *App) getSymbol(symbol string) (crypto.Symbol, bool) {
	a.Lock()
	defer a.Unlock()

	return a.getSymbolLocked(symbol)
}

func (a *App) getSymbolLocked(symbol string) (crypto.Symbol, bool) {
	for _, cd := range a.coinData {
		if coin, ok := cd.(*coin.CoinData); ok {
			if coin.Symbol.Symbol == symbol {
				return coin.Symbol, true
			}
		}
	}

	return crypto.Symbol{}, false
}

func (a *App) addSymbol(symbol crypto.Symbol) {
	var quotes []crypto.Quote
	if a.feed != nil {
		quotes, _ = a.feed.GetQuotes(a.currency, symbol.Symbol)
	}

	a.Lock()
	defer a.Unlock()

	if _, ok := a.getSymbolLocked(symbol.Symbol); ok {
		return
	}

	coin := coin.NewSymbol(symbol)
	if quotes != nil {
		coin = coin.UpdateQuote(quotes[0])
	}

	a.data.Append(coin)
}

func (a *App) updateQuote(quote crypto.Quote) {
	a.Lock()
	defer a.Unlock()
	for i, cd := range a.coinData {
		if coin, ok := cd.(*coin.CoinData); ok {
			if coin.Symbol == quote.Symbol {
				updatedCoin := coin.UpdateQuote(quote)
				if updatedCoin == nil {
					logger.Log.Error().Str("coin", coin.Symbol.Symbol).Str("quote", quote.Symbol.Symbol).Msg("Failed to update")
					continue
				}
				a.data.SetValue(i, updatedCoin)
				break
			}
		}
	}
}
