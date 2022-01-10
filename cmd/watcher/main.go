package main

import (
	"github.com/itohio/CoinWatcher/pkg/app"
)

func main() {
	watcher := app.New("Coin Watcher")
	watcher.Run()
}
