package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/CoinWatcher/pkg/widgets/coin"
)

func (a *App) makeList() *widget.List {
	list := widget.NewListWithData(a.data,
		func() fyne.CanvasObject {
			return coin.New(func(symbol string) {
				dialog.ShowConfirm(
					"Delete",
					fmt.Sprintf("You are about to delete %s.\nAre you sure?", symbol),
					func(b bool) {
						if b {
							a.delSymbol(symbol)
						}
					},
					a.window,
				)
			})
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*coin.CoinWidget).Bind(i.(binding.Untyped))
		},
	)

	return list
}
