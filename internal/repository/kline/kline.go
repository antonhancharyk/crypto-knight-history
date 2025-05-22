package kline

import (
	"fmt"
	"strings"
	"time"

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

func (t *Kline) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	var (
		klines []entity.Kline
		args   []any
		query  = "select * from klines where open_time >= $1 and open_time <= $2"
	)

	layout := "2006-01-02 15:04:05"

	var fromStr, toStr string

	if params.From == "" || params.To == "" {
		now := time.Now()
		firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-time.Second)

		fromStr = firstOfMonth.Format(layout)
		toStr = lastOfMonth.Format(layout)
	} else {
		fromStr = params.From
		toStr = params.To
	}

	fromTime, err := time.ParseInLocation(layout, fromStr, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("invalid 'from' datetime: %w", err)
	}
	toTime, err := time.ParseInLocation(layout, toStr, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("invalid 'to' datetime: %w", err)
	}

	fromMillis := fromTime.UnixMilli()
	toMillis := toTime.UnixMilli()

	args = append(args, fromMillis, toMillis)

	if params.Symbol != "" {
		query += " and symbol = $3"
		args = append(args, params.Symbol)
	}

	query += " order by open_time asc"

	err = t.db.Select(&klines, query, args...)

	return klines, err
}

func (t *Kline) GetLastKline() (entity.Kline, error) {
	klines := []entity.Kline{}

	err := t.db.Select(&klines, "select * from klines order by open_time desc limit 1")

	if len(klines) == 0 {
		return entity.Kline{}, err
	}

	return klines[0], err
}

func (k *Kline) CreateBulk(klines []entity.Kline) error {
	if len(klines) == 0 {
		return nil
	}

	var placeholders []string
	var values []any

	for i, kline := range klines {
		base := i * 12
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			base+1, base+2, base+3, base+4, base+5, base+6,
			base+7, base+8, base+9, base+10, base+11, base+12,
		))
		values = append(values,
			kline.Symbol,
			kline.OpenTime,
			kline.OpenPrice,
			kline.HighPrice,
			kline.LowPrice,
			kline.ClosePrice,
			kline.Volume,
			kline.CloseTime,
			kline.QuoteAssetVolume,
			kline.NumTrades,
			kline.TakerBuyBaseAssetVolume,
			kline.TakerBuyQuoteAssetVolume,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO klines (
			symbol, open_time, open_price, high_price, low_price,
			close_price, volume, close_time, quote_asset_volume,
			num_trades, taker_buy_base_asset_volume, taker_buy_quote_asset_volume
		)
		VALUES %s`, strings.Join(placeholders, ","))

	tx, err := k.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
