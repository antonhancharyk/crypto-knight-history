package kline

import (
	"database/sql"
	"errors"
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
		query  = "SELECT * FROM klines WHERE open_time >= $1 AND open_time <= $2 AND interval = $3"
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

	intervalDuration := intervalDuration(params.Interval)
	fromMillis := fromTime.Add(-480 * intervalDuration).UnixMilli()
	toMillis := toTime.UnixMilli()

	args = append(args, fromMillis, toMillis, params.Interval)

	if params.Symbol != "" {
		query += " AND symbol = $4"
		args = append(args, params.Symbol)
	}

	query += " ORDER BY open_time ASC, symbol ASC"

	err = t.db.Select(&klines, query, args...)

	return klines, err
}

func (t *Kline) GetLastKlineByInterval(interval string) (entity.Kline, error) {
	var k entity.Kline
	query := `
        SELECT * FROM klines
        WHERE interval = $1
        ORDER BY open_time DESC
        LIMIT 1
    `
	err := t.db.Get(&k, query, interval)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Kline{}, nil
		}
		return entity.Kline{}, err
	}

	return k, nil
}

func (k *Kline) CreateBulk(klines []entity.Kline) error {
	if len(klines) == 0 {
		return nil
	}

	var placeholders []string
	var values []any

	for i, kline := range klines {
		base := i * 13
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			base+1, base+2, base+3, base+4, base+5, base+6,
			base+7, base+8, base+9, base+10, base+11, base+12, base+13,
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
			kline.Interval,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO klines (
			symbol, open_time, open_price, high_price, low_price,
			close_price, volume, close_time, quote_asset_volume,
			num_trades, taker_buy_base_asset_volume, taker_buy_quote_asset_volume,
			interval
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

func intervalDuration(interval string) time.Duration {
	switch interval {
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return time.Hour
	}
}
