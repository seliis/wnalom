package main

import (
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/go-co-op/gocron"
)

var (
	Client      *futures.Client
	Config      Configuration
	MarketPrice float64
	Desc01m     *Description
	P           string
	L           int
	S           int
)

func main() {
	// Make Ready
	if err := ReadConfig(); err != nil {
		log.Println(err)
		return
	}

	P, L, S = "H", 0, 0

	// Using Test?
	futures.UseTestnet = Config.Test.Mode

	// Get Client
	if Config.Test.Mode {
		Client = futures.NewClient(Config.Test.Keys.API, Config.Test.Keys.Secret)
	} else {
		Client = futures.NewClient(Config.Keys.API, Config.Keys.Secret)
	}

	// Make Schedulers
	mainScheduler := gocron.NewScheduler(time.UTC)

	// Give Tasks to Schedulers
	mainScheduler.Every("1s").Do(func() {
		if marketPrice, err := GetMarketPrice(); err == nil {
			if converted, err := strconv.ParseFloat(marketPrice, 64); err == nil {
				MarketPrice = converted
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}

		if structPointer, err := GetDescription("1m"); err == nil {
			Desc01m = structPointer
		} else {
			panic(err)
		}

		RunStrategy()
	})

	// Run Schedulers
	mainScheduler.StartBlocking()
}
