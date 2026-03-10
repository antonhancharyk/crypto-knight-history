package kline

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	apperr "github.com/antonhancharyk/crypto-knight-history/internal/errors"
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/pkg/utilities"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dateLayout = "2006-01-02 15:04:05"
const lookbackIntervals = 480

type Kline struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Kline {
	return &Kline{db: db}
}

func (t *Kline) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return t.getKlinesByInterval(params, params.Interval)
}

func (t *Kline) GetKlines1h(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return t.getKlinesByInterval(params, "1h")
}

func (t *Kline) GetKlines30m(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return t.getKlinesByInterval(params, "30m")
}

func (t *Kline) getKlinesByInterval(params entity.GetKlinesQueryParams, interval string) ([]entity.Kline, error) {
	var klines []entity.Kline

	fromStr, toStr := resolveDateRange(params.From, params.To)
	fromTime, err := time.ParseInLocation(dateLayout, fromStr, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("invalid 'from' datetime: %w", apperr.ErrBadRequest)
	}
	toTime, err := time.ParseInLocation(dateLayout, toStr, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("invalid 'to' datetime: %w", apperr.ErrBadRequest)
	}

	intervalDur := utilities.IntervalDuration(interval)
	fromMillis := fromTime.Add(-lookbackIntervals * intervalDur).UnixMilli()
	toMillis := toTime.UnixMilli()

	query := "SELECT * FROM klines WHERE open_time >= $1 AND open_time <= $2 AND interval = $3"
	args := []any{fromMillis, toMillis, interval}

	if params.Symbol != "" {
		query += " AND symbol = $4"
		args = append(args, params.Symbol)
	}

	query += " ORDER BY open_time ASC, symbol ASC"

	err = t.db.Select(&klines, query, args...)
	return klines, err
}

func resolveDateRange(from, to string) (fromStr, toStr string) {
	if from != "" && to != "" {
		return from, to
	}
	now := time.Now().UTC()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	return firstOfMonth.Format(dateLayout), lastOfMonth.Format(dateLayout)
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
