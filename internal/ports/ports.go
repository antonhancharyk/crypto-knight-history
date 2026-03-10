package ports

import (
	"context"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	pbHistory "github.com/antongoncharik/crypto-knight-protos/gen/go/history"
	"google.golang.org/grpc"
)

// KlineRepository is the interface for kline data access.
type KlineRepository interface {
	GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetKlines1h(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetKlines30m(params entity.GetKlinesQueryParams) ([]entity.Kline, error)
	GetLastKlineByInterval(interval string) (entity.Kline, error)
	CreateBulk(klines []entity.Kline) error
}

// BinanceClient fetches data from Binance API.
type BinanceClient interface {
	Get(url string) ([]byte, error)
}

// HistoryClient calls the history gRPC service.
type HistoryClient interface {
	ProcessHistory(ctx context.Context, req *pbHistory.ProcessHistoryRequest, opts ...grpc.CallOption) (*pbHistory.ProcessHistoryResponse, error)
}

// HistoryProcessor runs history processing for given params (used by the task queue).
type HistoryProcessor interface {
	ProcessHistory(ctx context.Context, params entity.GetKlinesQueryParams) ([]entity.History, error)
}
