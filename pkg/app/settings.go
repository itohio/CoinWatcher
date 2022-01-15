package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"fyne.io/fyne/v2/storage"
	"github.com/itohio/CoinWatcher/pkg/crypto"
	"github.com/itohio/CoinWatcher/pkg/logger"
)

type Settings struct {
	APIKey   string        `json:"coinmarketcap_api_key"`
	Currency string        `json:"currency"`
	Interval time.Duration `json:"refresh_interval"`
}

func (a *App) defaultSettings() {
	logger.Log.Info().Msg("Loading default settings")
	a.apiKey = os.Getenv("COINWATCHER_KEY")
	a.feed = crypto.NewCMC(a.apiKey, a)
	a.currency = a.feed.GetCurrencies()[0]
	a.interval = time.Hour * 3

	a.saveSettings()
}

func (a *App) loadSettings() {
	reader, err := a.reader("config.json")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Could not get settings reader")
		a.defaultSettings()
		return
	}
	defer reader.Close()

	var settings Settings

	err = json.NewDecoder(reader).Decode(&settings)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Could not decode settings")
		a.defaultSettings()
		return
	}

	if settings.APIKey == "" {
		settings.APIKey = os.Getenv("COINWATCHER_KEY")
	}

	a.feed = crypto.NewCMC(settings.APIKey, a)
	a.currency = settings.Currency
	a.apiKey = settings.APIKey
	a.interval = settings.Interval
}

func (a *App) saveSettings() {
	settings := Settings{
		Currency: a.currency,
		Interval: a.interval,
		APIKey:   a.apiKey,
	}

	writer, err := a.writer("config.json")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Could not get settings writer")
		return
	}
	defer writer.Close()

	data, err := json.Marshal(&settings)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Could not marshal settings")
		return
	}

	if _, err := writer.Write(data); err != nil {
		logger.Log.Error().Err(err).Msg("Could not write settings")
	}
}

func (a *App) reader(base string) (io.ReadCloser, error) {
	uri, err := storage.Child(a.app.Storage().RootURI(), base)
	if err != nil {
		return nil, err
	}

	if ok, err := storage.CanRead(uri); !ok || err != nil {
		return nil, fmt.Errorf("Cannot read: %v", err)
	}

	reader, err := storage.Reader(uri)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (a *App) writer(base string) (io.WriteCloser, error) {
	uri, err := storage.Child(a.app.Storage().RootURI(), base)
	if err != nil {
		return nil, err
	}

	if ok, err := storage.CanWrite(uri); !ok || err != nil {
		return nil, fmt.Errorf("Cannot write: %v", err)
	}

	writer, err := storage.Writer(uri)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
