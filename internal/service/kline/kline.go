package kline

import (
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
)

type Kline struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Kline {
	return &Kline{repo: repo}
}

// func (t *Kline) GetAll(queryParams entity.QueryParams) ([]entity.Kline, error) {
// 	return t.repo.Kline.GetAll(queryParams)
// }

func (k *Kline) Create(kline entity.Kline) error {
	return k.repo.Kline.Create(kline)
}

func (k *Kline) CreateBulk(klines []entity.Kline) error {
	return k.repo.Kline.CreateBulk(klines)
}

// func (k *Kline) Get(sbl string) ([][]any, error) {
// 	params := url.Values{}
// 	params.Set("symbol", sbl)
// 	params.Set("interval", constant.INTERVAL_KLINES)
// 	params.Set("limit", constant.QUANTITY_KLINES)

// 	res, err := k.api.Get(constant.KLINES_URI+"?"+params.Encode(), false)
// 	if err != nil {
// 		return [][]any{}, err
// 	}

// 	data := [][]any{}
// 	err = json.Unmarshal(res, &data)
// 	if err != nil {
// 		return [][]any{}, err
// 	}

// 	return data, nil
// }
