package service

import (
	"context"
	"net/url"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/ports"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/kline"
	"github.com/antonhancharyk/crypto-knight-history/internal/service/queue"
)

type Kline interface {
	GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetBinanceKlines(params url.Values) ([]entity.Kline, error)
	CreateBulk(tracks []entity.Kline) error
	LoadInterval(interval string) error
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

func New(klineRepo ports.KlineRepository, binance ports.BinanceClient, history ports.HistoryClient) *Service {
	klineSvc := kline.New(klineRepo, binance, history)
	return &Service{
		Kline: klineSvc,
		Queue: queue.New(klineSvc),
	}
}

func (s *Service) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return s.Kline.GetKlines(params)
}

func (s *Service) CreateTask(params entity.GetKlinesQueryParams) *entity.Task {
	return s.Queue.CreateTask(params)
}

func (s *Service) GetTask(id string) (*entity.Task, bool) {
	return s.Queue.GetTask(id)
}
