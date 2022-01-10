package coin

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type table struct {
	sizes []float32
}

func (d *table) sumS() (w float32, c int) {
	for _, s := range d.sizes {
		if s >= 0 {
			w += s
		} else {
			c++
		}
	}
	return
}

func (d *table) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		if h < childSize.Height {
			h = childSize.Height
		}
	}

	w1, _ := d.sumS()
	if w < w1 {
		w = w1
	}

	return fyne.NewSize(w, h+theme.Padding())
}

func (d *table) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	sizes := make([]float32, len(objects))
	copy(sizes, d.sizes)
	w, c := d.sumS()
	dw := (containerSize.Width - w) / float32(c)
	if dw < 0 {
		dw = 1
	}
	for i, s := range sizes {
		if s < 0 {
			sizes[i] = dw
		}
	}

	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(sizes[i], 0))
	}
}
