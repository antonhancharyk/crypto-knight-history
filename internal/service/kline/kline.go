package kline

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	pbHistory "github.com/antongoncharik/crypto-knight-protos/gen/go/history"
	"github.com/antonhancharyk/crypto-knight-history/internal/constant"
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	grpcClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/grpc/client"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
	"github.com/antonhancharyk/crypto-knight-history/pkg/utilities"
)

type Kline struct {
	repo       *repository.Repository
	httpClient *client.HTTPClient
	grpcClient *grpcClient.Client
}

func New(repo *repository.Repository, httpClient *client.HTTPClient, grpcClient *grpcClient.Client) *Kline {
	return &Kline{repo: repo, httpClient: httpClient, grpcClient: grpcClient}
}

func (k *Kline) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return k.repo.Kline.GetKlines(params)
}

func (k *Kline) GetKlines1h(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return k.repo.Kline.GetKlines1h(params)
}

func (k *Kline) GetKlines30m(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return k.repo.Kline.GetKlines30m(params)
}

func (k *Kline) GetBinanceKlines(params url.Values) ([]entity.Kline, error) {
	sbl := params.Get("symbol")

	res, err := k.httpClient.Get(constant.KLINES_URI + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	data := [][]any{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	klinesBySymbol := []entity.Kline{}
	for _, v := range data {
		openTime := int64(v[0].(float64))
		openTimeUTC := utilities.FromUnixToUTC(openTime)
		openPrice, _ := strconv.ParseFloat(v[1].(string), 64)
		highPrice, _ := strconv.ParseFloat(v[2].(string), 64)
		lowPrice, _ := strconv.ParseFloat(v[3].(string), 64)
		closePrice, _ := strconv.ParseFloat(v[4].(string), 64)
		volume, _ := strconv.ParseFloat(v[5].(string), 64)
		closeTime := int64(v[6].(float64))
		closeTimeUTC := utilities.FromUnixToUTC(closeTime)
		quoteAssetVolume, _ := strconv.ParseFloat(v[7].(string), 64)
		numTrades := int64(v[8].(float64))
		takerBuyBaseAssetVolume, _ := strconv.ParseFloat(v[9].(string), 64)
		takerBuyQuoteAssetVolume, _ := strconv.ParseFloat(v[10].(string), 64)

		kline := entity.Kline{
			Symbol:                   sbl,
			OpenTime:                 openTime,
			OpenTimeUTC:              openTimeUTC,
			OpenPrice:                openPrice,
			HighPrice:                highPrice,
			LowPrice:                 lowPrice,
			ClosePrice:               closePrice,
			Volume:                   volume,
			CloseTime:                closeTime,
			CloseTimeUTC:             closeTimeUTC,
			QuoteAssetVolume:         quoteAssetVolume,
			NumTrades:                numTrades,
			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolume,
			TakerBuyQuoteAssetVolume: takerBuyQuoteAssetVolume,
		}

		klinesBySymbol = append(klinesBySymbol, kline)
	}

	return klinesBySymbol, nil
}

func (k *Kline) CreateBulk(klines []entity.Kline) error {
	return k.repo.Kline.CreateBulk(klines)
}

func (k *Kline) LoadInterval(interval string) error {
	last, err := k.repo.Kline.GetLastKlineByInterval(interval)
	if err != nil {
		return err
	}

	step := intervalDuration(interval)

	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if last.OpenTime != 0 {
		startTime = time.UnixMilli(last.OpenTime).Add(step).UTC()
	}

	now := time.Now().UTC()
	limitTime := now.Truncate(step)

	for startTime.Before(limitTime) {
		var wg sync.WaitGroup
		var mu sync.Mutex

		klines := make([][]entity.Kline, 0, len(constant.SYMBOLS))
		endTime := startTime.Add(498 * step)

		wg.Add(len(constant.SYMBOLS))
		for _, sbl := range constant.SYMBOLS {
			go func(symbol string) {
				defer wg.Done()

				params := url.Values{}
				params.Set("symbol", symbol)
				params.Set("interval", interval)
				params.Set("limit", "499")
				params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
				params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))

				data, err := k.GetBinanceKlines(params)
				if err != nil {
					log.Printf("[%s] %v", symbol, err)
					return
				}

				for i := range data {
					data[i].Interval = interval
				}

				if len(data) > 0 {
					mu.Lock()
					klines = append(klines, data)
					mu.Unlock()
				}

			}(sbl)
		}
		wg.Wait()

		if len(klines) == 0 {
			log.Printf("no klines for interval=%s starting from %s to %s", interval, startTime, endTime)
			startTime = startTime.Add(499 * step)
			time.Sleep(15 * time.Second)
			continue
		}

		wg.Add(len(klines))
		for _, arr := range klines {
			go func(kls []entity.Kline) {
				defer wg.Done()
				err := k.CreateBulk(kls)
				if err != nil {
					log.Print(err)
				}
			}(arr)
		}
		wg.Wait()

		if len(klines) > 0 && len(klines[0]) > 0 {
			log.Printf(
				"added interval=%s from %s to %s",
				interval,
				startTime.Format(time.RFC3339),
				endTime.Format(time.RFC3339),
			)
		}

		startTime = startTime.Add(499 * step)

		time.Sleep(15 * time.Second)
	}

	return nil
}

func (k *Kline) ProcessHistory(ctx context.Context, params entity.GetKlinesQueryParams) ([]entity.History, error) {
	var (
		histories []entity.History
		mu        sync.Mutex
	)

	klines30m, err := k.GetKlines30m(params)
	if err != nil {
		return nil, err
	}
	klines1h, err := k.GetKlines1h(params)
	if err != nil {
		return nil, err
	}

	results30m := make(map[string][]entity.Kline)
	for _, kline := range klines30m {
		results30m[kline.Symbol] = append(results30m[kline.Symbol], kline)
	}
	results1h := make(map[string][]entity.Kline)
	for _, kline := range klines1h {
		results1h[kline.Symbol] = append(results1h[kline.Symbol], kline)
	}

	sem := make(chan struct{}, 10)
	var wgGRPC sync.WaitGroup
	for resSymbol, resKlines := range results30m {
		if len(resKlines) == 0 {
			continue
		}

		wgGRPC.Add(1)
		go func(symbol string, klines []entity.Kline) {
			defer wgGRPC.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			var inputKlines []*pbHistory.InputKline
			for _, v := range klines {
				inputKlines = append(inputKlines, &pbHistory.InputKline{
					Id:                       v.Id,
					Symbol:                   v.Symbol,
					OpenTime:                 v.OpenTime,
					OpenPrice:                v.OpenPrice,
					HighPrice:                v.HighPrice,
					LowPrice:                 v.LowPrice,
					ClosePrice:               v.ClosePrice,
					Volume:                   v.Volume,
					CloseTime:                v.CloseTime,
					QuoteAssetVolume:         v.QuoteAssetVolume,
					NumTrades:                v.NumTrades,
					TakerBuyBaseAssetVolume:  v.TakerBuyBaseAssetVolume,
					TakerBuyQuoteAssetVolume: v.TakerBuyQuoteAssetVolume,
					Interval:                 v.Interval,
				})
			}

			kls := results1h[symbol]
			var inputKlines1h []*pbHistory.InputKline
			for _, v := range kls {
				inputKlines1h = append(inputKlines1h, &pbHistory.InputKline{
					Id:                       v.Id,
					Symbol:                   v.Symbol,
					OpenTime:                 v.OpenTime,
					OpenPrice:                v.OpenPrice,
					HighPrice:                v.HighPrice,
					LowPrice:                 v.LowPrice,
					ClosePrice:               v.ClosePrice,
					Volume:                   v.Volume,
					CloseTime:                v.CloseTime,
					QuoteAssetVolume:         v.QuoteAssetVolume,
					NumTrades:                v.NumTrades,
					TakerBuyBaseAssetVolume:  v.TakerBuyBaseAssetVolume,
					TakerBuyQuoteAssetVolume: v.TakerBuyQuoteAssetVolume,
					Interval:                 v.Interval,
				})
			}

			res, err := k.grpcClient.History.ProcessHistory(ctx, &pbHistory.ProcessHistoryRequest{
				Klines:   inputKlines,
				Klines1H: inputKlines1h,
			})
			if err != nil {
				fmt.Printf("failed to process history for %s: %v\n", symbol, err)
				return
			}

			h := entity.History{
				Symbol:                              res.Symbol,
				SumPositivePercentageChanges:        res.SumPositivePercentageChanges,
				CountPositiveChanges:                res.CountPositiveChanges,
				SumNegativePercentageChanges:        res.SumNegativePercentageChanges,
				CountNegativeChanges:                res.CountNegativeChanges,
				CountStopMarketOrders:               res.CountStopMarketOrders,
				CountIncomplitedStop:                res.CountIncomplitedStop,
				CountTransactions:                   res.CountTransactions,
				Grade:                               res.SumPositivePercentageChanges - res.SumNegativePercentageChanges - (float64(res.CountStopMarketOrders * 10)) - res.SumIncomplitedStopPercentageChanges,
				SumIncomplitedStopPercentageChanges: res.SumIncomplitedStopPercentageChanges,
			}

			mu.Lock()
			histories = append(histories, h)
			mu.Unlock()

		}(resSymbol, resKlines)
	}
	wgGRPC.Wait()

	var (
		sumPositivePercentageChanges        float64
		sumNegativePercentageChanges        float64
		sumIncomplitedStopPercentageChanges float64
		countPositiveChanges                int32
		countNegativeChanges                int32
		countStopMarketOrders               int32
		countIncomplitedStop                int32
		grade                               float64
	)

	for _, v := range histories {
		sumPositivePercentageChanges += v.SumPositivePercentageChanges
		sumNegativePercentageChanges += v.SumNegativePercentageChanges
		countPositiveChanges += v.CountPositiveChanges
		countNegativeChanges += v.CountNegativeChanges
		countStopMarketOrders += v.CountStopMarketOrders
		countIncomplitedStop += v.CountIncomplitedStop
		sumIncomplitedStopPercentageChanges += v.SumIncomplitedStopPercentageChanges
		grade += v.Grade
	}

	sort.Slice(histories, func(i, j int) bool {
		return histories[i].Grade < histories[j].Grade
	})

	h := entity.History{
		Symbol:                              "Total",
		SumPositivePercentageChanges:        sumPositivePercentageChanges,
		CountPositiveChanges:                countPositiveChanges,
		SumNegativePercentageChanges:        sumNegativePercentageChanges,
		CountNegativeChanges:                countNegativeChanges,
		CountStopMarketOrders:               countStopMarketOrders,
		SumIncomplitedStopPercentageChanges: sumIncomplitedStopPercentageChanges,
		CountIncomplitedStop:                countIncomplitedStop,
		CountTransactions:                   countPositiveChanges + countNegativeChanges,
		Grade:                               grade,
	}

	histories = append([]entity.History{h}, histories...)

	return histories, nil
}

func intervalDuration(interval string) time.Duration {
	switch interval {
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return time.Hour
	}
}
