package app

import (
	"fmt"
	"image"
	"image/png"
	"path"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/CoinWatcher/pkg/crypto"
	"github.com/itohio/CoinWatcher/pkg/logger"
)

type App struct {
	sync.Mutex
	app    fyne.App
	window fyne.Window

	currency    string
	interval    time.Duration
	lastUpdated time.Time
	feed        *crypto.Crypto
	apiKey      string

	coinData       []interface{}
	data           binding.ExternalUntypedList
	selectedSymbol string

	imageCache map[string]image.Image

	currencyWidget *widget.Select
	pbWidget       *widget.ProgressBar
	timeout        binding.Float
}

var _ crypto.Cache = &App{}

func New(name string) *App {
	a := app.NewWithID("CoinWatcher")
	w := a.NewWindow(name)
	w.Resize(fyne.NewSize(350, 600))

	ret := &App{
		app:        a,
		window:     w,
		imageCache: make(map[string]image.Image),
	}

	ret.loadSettings()

	ret.data = binding.BindUntypedList(&ret.coinData)
	ret.timeout = binding.NewFloat()

	ret.loadCoins()

	list := ret.makeList()
	menu := ret.makeMenu()

	ret.pbWidget = widget.NewProgressBarWithData(ret.timeout)
	ret.pbWidget.Min = 0
	ret.pbWidget.Max = 100
	ret.pbWidget.TextFormatter = func() string {
		return fmt.Sprint("ETA ", ((ret.interval-time.Since(ret.lastUpdated))/time.Minute)*time.Minute)
	}
	btnUpdate := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		ret.updateQuotes()
	})

	w.SetContent(
		container.NewBorder(
			menu,
			container.NewBorder(nil, nil, nil, btnUpdate, ret.pbWidget),
			nil, nil,
			list,
		),
	)

	return ret
}

func (a *App) Run() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				d := time.Since(a.lastUpdated).Seconds() / a.interval.Seconds()
				if d >= 1 {
					a.updateQuotes()
					d = 1
				}

				a.timeout.Set(d * 100)
			}
		}
	}()
	a.window.Show()
	a.app.Run()
}

func (a *App) LoadImage(url string) (image.Image, error) {
	base := fmt.Sprintf("cache_%s", path.Base(url))

	if img, ok := a.imageCache[base]; ok {
		return img, nil
	}

	reader, err := a.reader(base)
	if err != nil {
		logger.Log.Error().Msgf("Failed reading image from cache: %v", err)
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	a.imageCache[base] = img

	return img, nil
}

func (a *App) SaveImage(url string, img image.Image) {
	base := fmt.Sprintf("cache_%s", path.Base(url))
	a.imageCache[base] = img

	writer, err := a.writer(base)
	if err != nil {
		logger.Log.Error().Msgf("Failed writing image to cache: %v", err)
		return
	}
	defer writer.Close()

	err = png.Encode(writer, img)
	if err != nil {
		logger.Log.Error().Msgf("Failed encoding png: %v", err)
		return
	}
}
