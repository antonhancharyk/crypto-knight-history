package entity

import "time"

type Kline struct {
	Symbol                   string
	OpenTime                 int64
	OpenTimeUTC              time.Time
	OpenPrice                float64
	HighPrice                float64
	LowPrice                 float64
	ClosePrice               float64
	Volume                   float64
	CloseTime                int64
	CloseTimeUTC             time.Time
	QuoteAssetVolume         float64
	NumTrades                int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}
