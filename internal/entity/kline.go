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
	Interval                 string  `json:"interval" db:"interval"`
}

type MapKlines map[string][]Kline

type GetKlinesQueryParams struct {
	From     string `form:"from"`
	To       string `form:"to"`
	Interval string `form:"interval"`
	Symbol   string `form:"symbol,omitempty"`
}

type History struct {
	Symbol                              string  `json:"symbol"`
	SumPositivePercentageChanges        float64 `json:"sum_positive_percentage_changes"`
	CountPositiveChanges                int32   `json:"count_positive_changes"`
	SumNegativePercentageChanges        float64 `json:"sum_negative_percentage_changes"`
	CountNegativeChanges                int32   `json:"count_negative_changes"`
	CountStopMarketOrders               int32   `json:"count_stop_market_orders"`
	SumIncomplitedStopPercentageChanges float64 `json:"sum_incomplited_stop_percentage_changes"`
	CountIncomplitedStop                int32   `json:"count_incomplited_stop"`
	CountTransactions                   int32   `json:"count_transactions"`
	Grade                               float64 `json:"grade"`
}
