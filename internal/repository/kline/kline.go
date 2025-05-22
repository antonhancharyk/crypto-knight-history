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

	_, err := k.db.Exec(query, values...)

	return err
}
