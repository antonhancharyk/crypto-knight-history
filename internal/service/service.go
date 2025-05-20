package service

import (
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/kline"
)

type Kline interface {
	// GetAll(queryParams entity.QueryParams) ([]entity.Kline, error)
	Create(track entity.Kline) error
	CreateBulk(tracks []entity.Kline) error
}

type Service struct {
	Kline
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Kline: kline.New(repo),
	}
}
