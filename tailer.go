package main

import (
	"context"
	"strconv"

	"github.com/cinar/indicator"
)

func GetMarketPrice() (string, error) {
	if data, err := Client.NewListPricesService().Symbol(Config.Symbol).Do(context.Background()); err == nil {
		return data[0].Price, nil
	} else {
		return "0", err
	}
}

type Description struct {
	Volume []int64
	Close  []float64
	High   []float64
	Low    []float64
	CCI    []float64
	RSI    []float64
	MFI    []float64
	KDJ    struct {
		K []float64
		D []float64
		J []float64
	}
	MACD struct {
		MACD   []float64
		Signal []float64
	}
	BollingerBands struct {
		Middle []float64
		Upper  []float64
		Lower  []float64
	}
}

func GetDescription(interval string) (*Description, error) {
	desc := new(Description)

	convert := func(str string) float64 {
		res, _ := strconv.ParseFloat(str, 64)
		return res
	}

	if data, err := Client.NewKlinesService().Symbol(Config.Symbol).Interval(interval).Do(context.Background()); err == nil {
		for i := 0; i < len(data); i++ {
			desc.Volume = append(desc.Volume, int64(convert(data[i].Volume)))
			desc.Close = append(desc.Close, convert(data[i].Close))
			desc.High = append(desc.High, convert(data[i].High))
			desc.Low = append(desc.Low, convert(data[i].Low))
		}

		// CCI
		desc.CCI = indicator.CommunityChannelIndex(9, desc.High, desc.Low, desc.Close)

		// RSI
		_, desc.RSI = indicator.Rsi(desc.Close)

		// MFI
		desc.MFI = indicator.MoneyFlowIndex(14, desc.High, desc.Low, desc.Close, desc.Volume)

		// KDJ
		desc.KDJ.K, desc.KDJ.D, desc.KDJ.J = indicator.DefaultKdj(desc.High, desc.Low, desc.Close)

		// MACD
		desc.MACD.MACD, desc.MACD.Signal = indicator.Macd(desc.Close)

		// Bollinger Bands
		desc.BollingerBands.Middle, desc.BollingerBands.Upper, desc.BollingerBands.Lower = indicator.BollingerBands(desc.Close)

		return desc, nil
	} else {
		return desc, err
	}
}
