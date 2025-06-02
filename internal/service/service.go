package service

import (
	"context"
	"net/url"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/kline"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/queue"
	grpcClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/grpc/client"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
)

type Kline interface {
	GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetLastKline() (entity.Kline, error)
	GetBinanceKlines(params url.Values) ([]entity.Kline, error)
	CreateBulk(tracks []entity.Kline) error
	LoadKlinesForPeriod() error
	ProcessHistory(ctx context.Context, params entity.GetKlinesQueryParams) ([]entity.History, error)
}

type Queue interface {
	CreateTask(params entity.GetKlinesQueryParams) *entity.Task
	GetTask(id string) (*entity.Task, bool)
}

type Service struct {
	Kline
	Queue
}

func New(repo *repository.Repository, httpClient *client.HTTPClient, grpcClient *grpcClient.Client) *Service {
	klineSvc := kline.New(repo, httpClient, grpcClient)

	return &Service{
		Kline: klineSvc,
		Queue: queue.New(klineSvc),
	}
}
