package kline

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/constant"
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
	"github.com/antonhancharyk/crypto-knight-history/pkg/utilities"
)

type Kline struct {
	repo       *repository.Repository
	httpClient *client.HTTPClient
}

func New(repo *repository.Repository, httpClient *client.HTTPClient) *Kline {
	return &Kline{repo: repo, httpClient: httpClient}
}

func (k *Kline) GetKlines(params entity.GetKlinesQueryParams) ([]entity.Kline, error) {
	return k.repo.Kline.GetKlines(params)
}

func (k *Kline) GetLastKline() (entity.Kline, error) {
	return k.repo.Kline.GetLastKline()
}

func (k *Kline) GetBinanceKlines(params url.Values) []entity.Kline {
	sbl := params.Get("symbol")

	res, err := k.httpClient.Get(constant.KLINES_URI + "?" + params.Encode())
	if err != nil {
		log.Fatal(err)
	}

	data := [][]any{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Fatal(err)
	}

	lstIdx := len(data) - 1
	klinesBySymbol := []entity.Kline{}
	for k, v := range data {
		if k == lstIdx {
			break
		}

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

	return klinesBySymbol
}

func (k *Kline) CreateBulk(klines []entity.Kline) error {
	return k.repo.Kline.CreateBulk(klines)
}

func (k *Kline) LoadKlinesForPeriod() {
	lastKline, err := k.GetLastKline()
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	if lastKline.OpenTime != 0 {
		log.Printf("find last kline [unix open time]: %d", lastKline.OpenTime)
		startTime = time.UnixMilli(lastKline.OpenTime).Add(1 * time.Hour).UTC()
	} else {
		log.Printf("last kline is empty [unix start time]: %d", startTime.UnixMilli())
	}
	now := time.Now().UTC()
	endTime := now.Truncate(time.Hour)

	var wg sync.WaitGroup

	for startTime.Before(endTime) {
		klines := [][]entity.Kline{}
		wg.Add(len(constant.SYMBOLS))
		for _, symbol := range constant.SYMBOLS {
			go func(sbl string) {
				defer wg.Done()

				params := url.Values{}
				params.Set("symbol", symbol)
				params.Set("interval", constant.INTERVAL_KLINES)
				params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
				params.Set("endTime", strconv.FormatInt(startTime.Add(20*24*time.Hour).UnixMilli(), 10))

				klns := k.GetBinanceKlines(params)

				klines = append(klines, klns)
			}(symbol)
		}
		wg.Wait()

		wg.Add(len(klines))
		for _, v := range klines {
			go func(kl []entity.Kline) {
				defer wg.Done()

				k.CreateBulk(kl)
			}(v)
		}
		wg.Wait()

		if len(klines) != 0 {
			log.Printf("add new klines [quantity]: %d, [unix open time from]: %d, [unix open time to]: %d", len(klines[0]), klines[0][0].OpenTime, klines[0][len(klines[0])-1].OpenTime)
		}

		startTime = startTime.Add(20 * 24 * time.Hour)
		time.Sleep(15 * time.Second)
	}
}
