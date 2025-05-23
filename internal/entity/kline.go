package entity

import "time"

type Kline struct {
	Id                       int64     `json:"id" db:"id"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	Symbol                   string    `json:"symbol" db:"symbol"`
	OpenTime                 int64     `json:"open_time" db:"open_time"`
	OpenTimeUTC              time.Time
	OpenPrice                float64 `json:"open_price" db:"open_price"`
	HighPrice                float64 `json:"high_price" db:"high_price"`
	LowPrice                 float64 `json:"low_price" db:"low_price"`
	ClosePrice               float64 `json:"close_price" db:"close_price"`
	Volume                   float64 `json:"volume" db:"volume"`
	CloseTime                int64   `json:"close_time" db:"close_time"`
	CloseTimeUTC             time.Time
	QuoteAssetVolume         float64 `json:"quote_asset_volume" db:"quote_asset_volume"`
	NumTrades                int64   `json:"num_trades" db:"num_trades"`
	TakerBuyBaseAssetVolume  float64 `json:"taker_buy_base_asset_volume" db:"taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume float64 `json:"taker_buy_quote_asset_volume" db:"taker_buy_quote_asset_volume"`
}

type MapKlines map[string][]Kline

type GetKlinesQueryParams struct {
	From   string `form:"from"`
	To     string `form:"to"`
	Symbol string `form:"symbol"`
}

type History struct {
	Symbol                    string  `json:"symbol"`
	AmountPositivePercentages float64 `json:"amount_positive_percentages"`
	AmountNegativePercentages float64 `json:"amount_negative_percentages"`
	QuantityStopMarkets       float64 `json:"quantity_stop_markets"`
}
