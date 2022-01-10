# CoinWatcher
A simple GUI to track cryptocurrency prices.

![Typical Coin Watcher GUI](media/watcher.jpg)
# How to run it

```
$ go install github.com/itohi/CoinWatcher/cmd/watcher@latest
```

All the requirements for building Fyne apps apply.

## Configuration

You will need Coinmarketcap API key for this to work. They offer a free tier with 
10000 API credits per month, so it should be plenty enough for casual use.

You can setup the API key using `COINWATCHER_KEY` environment variable at first start. Otherwise it is possible
to configure the api key using settings button.

