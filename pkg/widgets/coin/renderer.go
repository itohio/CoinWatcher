package coin

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const iconSize = 24

type coinRenderer struct {
	symbol    *canvas.Text
	name      *canvas.Text
	price     *canvas.Text
	marketCap *canvas.Text
	volume    *canvas.Text
	pc1H      *canvas.Text
	pc24H     *canvas.Text
	pc7D      *canvas.Text
	pc30D     *canvas.Text
	icon      *canvas.Image

	widget    *CoinWidget
	container *fyne.Container
}

func (w *CoinWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	symbol := canvas.NewText(w.symbol, theme.ForegroundColor())
	name := canvas.NewText(w.name, theme.ForegroundColor())
	price := canvas.NewText("", theme.ForegroundColor())
	price.Alignment = fyne.TextAlignTrailing

	volume := canvas.NewText("", theme.ForegroundColor())
	volume.Alignment = fyne.TextAlignTrailing
	marketCap := canvas.NewText("", theme.ForegroundColor())
	marketCap.Alignment = fyne.TextAlignTrailing

	pc1H := canvas.NewText("", theme.ForegroundColor())
	pc24H := canvas.NewText("", theme.ForegroundColor())
	pc7D := canvas.NewText("", theme.ForegroundColor())
	pc30D := canvas.NewText("", theme.ForegroundColor())

	ret := &coinRenderer{
		widget:    w,
		symbol:    symbol,
		name:      name,
		price:     price,
		volume:    volume,
		marketCap: marketCap,
		pc1H:      pc1H,
		pc24H:     pc24H,
		pc7D:      pc7D,
		pc30D:     pc30D,
	}

	ret.refreshNumbers()
	ret.updateObjects()

	return ret
}

func (r *coinRenderer) updateObjects() {
	var icon fyne.CanvasObject = r.widget.icon
	symbol := container.NewVBox(r.symbol, r.name)
	price := container.NewVBox(r.price, r.volume, r.marketCap)

	var objs []fyne.CanvasObject

	if r.widget.icon != nil {
		objs = append(objs, container.NewPadded(icon))
	}
	if r.widget.onMenu != nil {
		btn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			if r.widget.onMenu != nil {
				r.widget.onMenu(r.widget.symbol)
			}
		})
		btn.Importance = widget.LowImportance
		objs = append(objs, btn)
	}

	name := container.NewHBox(container.NewVBox(objs...), symbol)
	stats := container.NewVBox(r.pc1H, r.pc24H, r.pc7D, r.pc30D)

	a := canvas.NewText("493.9999K", theme.ForegroundColor()).MinSize()
	b := canvas.NewText("W: 5555.5", theme.ForegroundColor()).MinSize()

	root := container.New(&table{sizes: []float32{-1, a.Width, b.Width}}, name, price, stats)
	r.container = root

	r.applyTheme()
}

func (r *coinRenderer) Layout(size fyne.Size) {
	r.container.Layout.Layout(r.container.Objects, size)
}

func (r *coinRenderer) MinSize() fyne.Size {
	size := r.container.MinSize()
	if size.Width < 10*iconSize {
		size.Width = 10 * iconSize
	}
	size.Height = 84
	return size
}

func formatNumber(prefix string, n float64, accuracy int) string {
	format := fmt.Sprintf("%%s%%0.%df%%s", accuracy)
	switch {
	case n >= 1e12:
		return fmt.Sprintf(format, prefix, n/1e12, "T")
	case n >= 1e9:
		return fmt.Sprintf(format, prefix, n/1e9, "B")
	case n >= 1e6:
		return fmt.Sprintf(format, prefix, n/1e6, "M")
	case n >= 1e3:
		return fmt.Sprintf(format, prefix, n/1e3, "K")
	case n < 1e-3:
		return fmt.Sprintf(format, prefix, n*1e6, "u")
	case n < 1:
		return fmt.Sprintf(format, prefix, n*1e3, "m")
	default:
		return fmt.Sprintf(format, prefix, n, "")
	}
}

func (r *coinRenderer) refreshNumbers() {
	r.symbol.Text = r.widget.symbol
	r.name.Text = r.widget.name
	r.price.Text = formatNumber("", r.widget.price, 2)
	r.marketCap.Text = formatNumber("mc: ", r.widget.marketCap, 1)
	r.volume.Text = formatNumber("V: ", r.widget.volume, 1)

	r.pc1H.Text = fmt.Sprintf("H: %0.1f", r.widget.pc1H)
	r.pc24H.Text = fmt.Sprintf("D: %0.1f", r.widget.pc24H)
	r.pc7D.Text = fmt.Sprintf("W: %0.1f", r.widget.pc7D)
	r.pc30D.Text = fmt.Sprintf("M: %0.1f", r.widget.pc30D)
}

func (r *coinRenderer) Refresh() {
	r.refreshNumbers()
	r.updateObjects()
	r.container.Refresh()

	r.Layout(r.widget.Size())
	canvas.Refresh(r.widget)
}

func (r *coinRenderer) Objects() []fyne.CanvasObject {
	return r.container.Objects
}

func (r *coinRenderer) Destroy() {

}

func (r *coinRenderer) applyTheme() {

	r.symbol.TextSize = theme.TextHeadingSize()
	r.symbol.Color = theme.ForegroundColor()

	r.name.TextSize = theme.TextSubHeadingSize()
	r.name.Color = theme.ForegroundColor()

	r.price.TextSize = theme.TextSize()
	r.volume.TextSize = theme.TextSize() * 2.0 / 3.0
	r.applyThemeChange(r.price, r.widget.pc1H)
	r.applyThemeChange(r.volume, r.widget.pc24H)

	r.marketCap.TextSize = theme.TextSize() * 2.0 / 3.0
	r.pc1H.TextSize = theme.TextSize() * 2.0 / 3.0
	r.pc24H.TextSize = theme.TextSize() * 2.0 / 3.0
	r.pc7D.TextSize = theme.TextSize() * 2.0 / 3.0
	r.pc30D.TextSize = theme.TextSize() * 2.0 / 3.0

	r.applyThemeChange(r.pc1H, r.widget.pc1H)
	r.applyThemeChange(r.pc24H, r.widget.pc24H)
	r.applyThemeChange(r.pc7D, r.widget.pc7D)
	r.applyThemeChange(r.pc30D, r.widget.pc30D)
}

func (r *coinRenderer) applyThemeChange(obj *canvas.Text, change float64) {
	green := color.RGBA{0, 255, 0, 0}
	red := color.RGBA{255, 0, 0, 0}

	if change > 0 {
		obj.Color = green
	} else if change < 0 {
		obj.Color = red
	} else {
		obj.Color = theme.ForegroundColor()
	}
}
