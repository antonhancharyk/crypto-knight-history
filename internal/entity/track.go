package entity

import (
	"time"

	"github.com/lib/pq"
)

type Track struct {
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	Symbol     string          `json:"symbol" db:"symbol" validate:"required"`
	HighPrice  float64         `json:"high_price" db:"high_price" validate:"numeric"`
	LowPrice   float64         `json:"low_price" db:"low_price" validate:"numeric"`
	HighPrices pq.Float64Array `json:"high_prices" db:"high_prices"`
	LowPrices  pq.Float64Array `json:"low_prices" db:"low_prices"`
}
