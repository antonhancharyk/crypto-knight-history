package kline

import (
	"encoding/json"
	"errors"
	"net/url"
	"testing"

	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
)

type mockKlineRepo struct {
	klines []entity.Kline
	err    error
}

func (m *mockKlineRepo) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error)     { return m.klines, m.err }
func (m *mockKlineRepo) GetKlines1h(params entity.GetKlinesQueryParams) ([]entity.Kline, error)   { return m.klines, m.err }
func (m *mockKlineRepo) GetKlines30m(params entity.GetKlinesQueryParams) ([]entity.Kline, error) { return m.klines, m.err }
func (m *mockKlineRepo) GetLastKlineByInterval(interval string) (entity.Kline, error)            { return entity.Kline{}, m.err }
func (m *mockKlineRepo) CreateBulk(klines []entity.Kline) error                                  { return m.err }

type mockBinanceClient struct {
	body []byte
	err  error
}

func (m *mockBinanceClient) Get(url string) ([]byte, error) {
	return m.body, m.err
}

func TestGetKlines_DelegatesToRepo(t *testing.T) {
	want := []entity.Kline{{Symbol: "BTCUSDT", Interval: "1h"}}
	repo := &mockKlineRepo{klines: want}
	svc := New(repo, &mockBinanceClient{}, nil)

	got, err := svc.GetKlines(entity.GetKlinesQueryParams{Interval: "1h"})
	if err != nil {
		t.Fatalf("GetKlines: %v", err)
	}
	if len(got) != 1 || got[0].Symbol != want[0].Symbol {
		t.Errorf("GetKlines = %v, want %v", got, want)
	}
}

func TestGetKlines_ReturnsRepoError(t *testing.T) {
	repoErr := errors.New("repo error")
	repo := &mockKlineRepo{err: repoErr}
	svc := New(repo, &mockBinanceClient{}, nil)
	_, err := svc.GetKlines(entity.GetKlinesQueryParams{})
	if err != repoErr {
		t.Errorf("err = %v", err)
	}
}

// Binance kline row: [openTime, open, high, low, close, volume, closeTime, quoteVol, numTrades, takerBuyBase, takerBuyQuote]
var binanceKlineRow = []any{
	float64(1609459200000), "50000", "51000", "49000", "50500", "1000",
	float64(1609459259999), "50000000", int64(100), "500", "500",
}

func TestGetBinanceKlines_ParsesResponse(t *testing.T) {
	body := mustMarshal([][]any{binanceKlineRow})
	binance := &mockBinanceClient{body: body}
	svc := New(&mockKlineRepo{}, binance, nil)

	params := url.Values{}
	params.Set("symbol", "BTCUSDT")
	klines, err := svc.GetBinanceKlines(params)
	if err != nil {
		t.Fatalf("GetBinanceKlines: %v", err)
	}
	if len(klines) != 1 {
		t.Fatalf("len(klines) = %d", len(klines))
	}
	k := klines[0]
	if k.Symbol != "BTCUSDT" {
		t.Errorf("Symbol = %q", k.Symbol)
	}
	if k.OpenPrice != 50000 || k.HighPrice != 51000 || k.ClosePrice != 50500 {
		t.Errorf("prices: open=%.0f high=%.0f close=%.0f", k.OpenPrice, k.HighPrice, k.ClosePrice)
	}
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
