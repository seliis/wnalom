package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"
	"github.com/go-co-op/gocron"
)

var (
	Client *futures.Client
)

func main() {
	// Make Ready
	if err := ReadConfig(); err != nil {
		return
	}

	// Is For Test Or Else?
	if Config.Test.Mode {
		Client = futures.NewClient(Config.Test.Keys.API, Config.Test.Keys.Secret)
		futures.UseTestnet = true
	} else {
		Client = futures.NewClient(Config.Keys.API, Config.Keys.Secret)
		futures.UseTestnet = false
	}

	// Make Schedulers
	mainScheduler := gocron.NewScheduler(time.UTC)

	// Give Tasks to Schedulers
	mainScheduler.Every("1s").Do(func() {
		if data, err := Client.NewListPricesService().Symbol(Config.Symbol).Do(context.Background()); err == nil {
			log.Println("Now Price:", data[0].Price)
		} else {
			log.Println("Can't Get Price of", Config.Symbol)
		}

		if data, err := Client.NewKlinesService().Symbol(Config.Symbol).Interval("1m").Do(context.Background()); err == nil {
			dict := map[string][]float64{
				"CLOSE": {},
				"HIGH":  {},
				"LOW":   {},
			}
			for i := 0; i < len(data); i++ {
				close, _ := strconv.ParseFloat(data[i].Close, 64)
				high, _ := strconv.ParseFloat(data[i].High, 64)
				low, _ := strconv.ParseFloat(data[i].Low, 64)
				dict["CLOSE"] = append(dict["CLOSE"], close)
				dict["HIGH"] = append(dict["HIGH"], high)
				dict["LOW"] = append(dict["LOW"], low)
			}

			// RSI
			_, rsi := indicator.Rsi(dict["CLOSE"])
			fmt.Println("RSI:", int(rsi[len(rsi)-1]))

			// KDJ
			k, d, j := indicator.DefaultKdj(dict["HIGH"], dict["LOW"], dict["CLOSE"])
			fmt.Println(
				"K:", int(k[len(k)-1]),
				"D:", int(d[len(d)-1]),
				"J:", int(j[len(j)-1]),
			)

			// MACD
			macd, signal := indicator.Macd(dict["CLOSE"])
			fmt.Println(
				"MACD:", int(macd[len(macd)-1]*10000),
				"SIGNAL:", int(signal[len(signal)-1]*10000),
			)
		} else {
			log.Println("Can't Get Klines of", Config.Symbol)
		}
	})

	// Run Schedulers
	mainScheduler.StartBlocking()
}
