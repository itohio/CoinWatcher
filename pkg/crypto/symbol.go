package crypto

import (
	"image"
	"net/http"
	"time"

	"github.com/itohio/CoinWatcher/pkg/logger"
)

type Symbol struct {
	iconCache Cache
	Id        int
	Name      string
	Symbol    string
	IconURL   string
}

type Quote struct {
	Symbol

	Price            float64
	Volume24H        float64
	Volume7D         float64
	Volume30D        float64
	Volume24Hbase    float64
	Volume24Hquote   float64
	PercentChange1H  float64
	PercentChange24H float64
	PercentChange7D  float64
	PercentChange30D float64
	MarketCap        float64
	LastUpdated      time.Time
}

type Ohlcv struct {
	Symbol

	LastUpdated string
	TimeOpen    string
	TimeClose   string
	Quote       map[string]OhlcvQuote
}

type OhlcvQuote struct {
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	Timestamp   time.Time
	LastUpdated time.Time
}

func (s *Symbol) Icon() image.Image {
	if s.iconCache != nil {
		if img, err := s.iconCache.LoadImage(s.IconURL); err == nil {
			return img
		}
	}

	response, err := http.Get(s.IconURL)
	if err != nil {
		logger.Log.Warn().Str("url", s.IconURL).Err(err).Msg("No icon")
		return image.NewGray(image.Rectangle{Max: image.Point{32, 32}})
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		logger.Log.Warn().Str("url", s.IconURL).Err(err).Msg("No icon")
		return image.NewGray(image.Rectangle{Max: image.Point{32, 32}})
	}

	logger.Log.Debug().Str("url", s.IconURL).Float32("w", float32(img.Bounds().Dx())).Float32("h", float32(img.Bounds().Dy())).Msg("icon img")

	if s.iconCache != nil {
		s.iconCache.SaveImage(s.IconURL, img)
	}

	return img
}
