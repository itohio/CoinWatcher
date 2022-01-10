package coin

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CoinWidget struct {
	widget.BaseWidget
	sync.Mutex

	iconUrl   string
	icon      *canvas.Image
	symbol    string
	name      string
	price     float64
	volume    float64
	pc1H      float64
	pc24H     float64
	pc7D      float64
	pc30D     float64
	marketCap float64

	data   binding.DataItem
	onMenu func(string)

	showStats bool
}

func New(onMenu func(string)) *CoinWidget {
	ret := &CoinWidget{
		onMenu: onMenu,
		icon:   canvas.NewImageFromResource(theme.FileImageIcon()),
	}
	ret.ExtendBaseWidget(ret)

	return ret
}

// MinSize returns the size that this widget should not shrink below.
//
// Implements: fyne.Widget
func (w *CoinWidget) MinSize() fyne.Size {
	w.ExtendBaseWidget(w)
	return w.BaseWidget.MinSize()
}

func (w *CoinWidget) Tapped(*fyne.PointEvent) {
	w.showStats = !w.showStats
	w.Refresh()
}
