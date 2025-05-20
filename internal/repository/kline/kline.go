package kline

import (
	"fmt"
	"strings"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Kline struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Kline {
	return &Kline{db: db}
}

// func (t *Kline) GetAll(queryParams entity.QueryParams) ([]entity.Kline, error) {
// 	tracksData := []entity.Kline{}

// 	if queryParams.Full && queryParams.Symbol != "" {
// 		err := t.db.Select(&tracksData, "select symbol, high_price, low_price, created_at, is_order, high_created_at, low_created_at, high_prices, low_prices, take_profit_high_prices, take_profit_low_prices from klines where created_at between $1 and $2 and symbol = $3 order by created_at desc", queryParams.From, queryParams.To, queryParams.Symbol)

// 		return tracksData, err
// 	}

// 	if queryParams.Full && queryParams.Symbol == "" {
// 		err := t.db.Select(&tracksData, "select symbol, high_price, low_price, created_at, is_order, high_created_at, low_created_at, high_prices, low_prices, take_profit_high_prices, take_profit_low_prices from tracks where created_at between $1 and $2 order by created_at desc", queryParams.From, queryParams.To)

// 		return tracksData, err
// 	}

// 	if queryParams.Symbol != "" {
// 		err := t.db.Select(&tracksData, "with ranked_tracks as (select id, symbol, high_price, low_price, created_at, lag(high_price) over (partition by symbol order by created_at) as prev_high_price, lag(low_price) over (partition by symbol order by created_at) as prev_low_price, high_created_at, low_created_at, high_prices, low_prices, take_profit_high_prices, take_profit_low_prices from tracks where created_at between $1 and $2 and symbol = $3) select symbol, high_price, low_price, created_at, high_created_at, low_created_at, high_prices, low_prices, take_profit_high_prices, take_profit_low_prices from ranked_tracks where high_price != prev_high_price or low_price != prev_low_price or prev_high_price is null order by created_at desc", queryParams.From, queryParams.To, queryParams.Symbol)

// 		return tracksData, err
// 	}

// 	err := t.db.Select(&tracksData, "with ranked_tracks as (select id, symbol, high_price, low_price, created_at, lag(high_price) over (partition by symbol order by created_at) as prev_high_price, lag(low_price) over (partition by symbol order by created_at) as prev_low_price, high_created_at, low_created_at, high_prices, low_prices, take_profit_high_prices, take_profit_low_prices from tracks where created_at between $1 and $2) select symbol, high_price, low_price, created_at, high_created_at, low_created_at, high_prices, low_prices, take_profit_high_prices, take_profit_low_prices from ranked_tracks where high_price != prev_high_price or low_price != prev_low_price or prev_high_price is null order by created_at desc", queryParams.From, queryParams.To)

// 	return tracksData, err
// }

func (k *Kline) Create(track entity.Kline) error {
	var err error

	_, err = k.db.Exec(
		`INSERT INTO klines (symbol)
		VALUES ($1)`, track.Symbol)

	return err
}

func (k *Kline) CreateBulk(klines []entity.Kline) error {
	var placeholders []string
	var values []any
	for i, kline := range klines {
		placeholders = append(placeholders, fmt.Sprintf("($%d)", i*11+1))
		values = append(values, kline.Symbol)
	}

	var err error
	query := fmt.Sprintf(
		`INSERT INTO klines (symbol)
		VALUES %s`, strings.Join(placeholders, ","))

	_, err = k.db.Exec(query, values...)

	return err
}
