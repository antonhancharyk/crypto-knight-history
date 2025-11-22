package repository

import (
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository/kline"
	"github.com/jmoiron/sqlx"
)

type Kline interface {
	GetKlines(getKlinesQueryParams entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetLastKlineByInterval(interval string) (entity.Kline, error)
	CreateBulk(klines []entity.Kline) error
}

type Repository struct {
	Kline
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Kline: kline.New(db),
	}
}
