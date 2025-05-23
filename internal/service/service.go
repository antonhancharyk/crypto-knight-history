package service

import (
	"context"
	"net/url"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/kline"
	grpcClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/grpc/client"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
)

type Kline interface {
	GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetLastKline() (entity.Kline, error)
	GetBinanceKlines(params url.Values) []entity.Kline
	CreateBulk(tracks []entity.Kline) error
	LoadKlinesForPeriod()
	ProcessHistory(ctx context.Context) (*entity.History, error)
}

type Service struct {
	Kline
}

func New(repo *repository.Repository, httpClient *client.HTTPClient, grpcClient *grpcClient.Client) *Service {
	return &Service{
		Kline: kline.New(repo, httpClient, grpcClient),
	}
}
