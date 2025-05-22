package service

import (
	"net/url"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/kline"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
)

type Kline interface {
	GetKlines(params url.Values) []entity.Kline
	GetLastKline() (entity.Kline, error)
	CreateBulk(tracks []entity.Kline) error
	LoadKlinesForPeriod()
}

type Service struct {
	Kline
}

func New(repo *repository.Repository, httpClient *client.HTTPClient) *Service {
	return &Service{
		Kline: kline.New(repo, httpClient),
	}
}
