package coin

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/CoinWatcher/pkg/crypto"
)

var ErrBadData = errors.New("bad data")

type CoinData struct {
	crypto.Symbol
	crypto.Quote
}

func NewSymbol(symbol crypto.Symbol) *CoinData {
	return &CoinData{
		Symbol: symbol,
	}
}

func (w *CoinWidget) Bind(data binding.Untyped) {
	w.Lock()
	oldData := w.data
	w.data = data
	w.Unlock()
	if oldData != nil {
		oldData.RemoveListener(w)
	}
	data.AddListener(w)
}

func (w *CoinWidget) getData() (*CoinData, error) {
	w.Lock()
	defer w.Unlock()

	if w.data == nil {
		return nil, ErrBadData
	}
	source, ok := w.data.(binding.Untyped)
	if !ok {
		return nil, ErrBadData
	}
	dataSource, err := source.Get()
	if !ok || err != nil {
		return nil, fmt.Errorf("bad data: %v", err)
	}
	data, ok := dataSource.(*CoinData)
	if !ok || data == nil {
		return nil, ErrBadData
	}

	return data, nil
}

func (w *CoinWidget) DataChanged() {
	data, err := w.getData()
	if err != nil {
		return
	}
	if w.iconUrl != data.IconURL || w.icon == nil {
		w.iconUrl = data.IconURL
		w.icon = canvas.NewImageFromImage(data.Icon())
		w.icon.SetMinSize(fyne.NewSize(32, 32))
		//w.icon.Resize(fyne.NewSize(32, 32))
		w.icon.FillMode = canvas.ImageFillContain
	}
	w.symbol = data.Symbol.Symbol
	w.name = data.Symbol.Name
	w.price = data.Price
	w.volume = data.Volume24H
	w.marketCap = data.MarketCap
	w.pc1H = data.PercentChange1H
	w.pc24H = data.PercentChange24H
	w.pc7D = data.PercentChange7D
	w.pc30D = data.PercentChange30D
	w.Refresh()
}

func (c *CoinData) UpdateQuote(q crypto.Quote) *CoinData {
	if c.Symbol.Symbol != q.Symbol.Symbol {
		return nil
	}
	return &CoinData{
		Symbol: q.Symbol,
		Quote:  q,
	}
}
