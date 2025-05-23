package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/antonhancharyk/crypto-knight-history/internal/db"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service"
	grpcClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/grpc/client"
	httpClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
	"github.com/antonhancharyk/crypto-knight-history/pkg/utilities"
	"github.com/joho/godotenv"
)

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("History is running...")

	godotenv.Load()

	dbClient := db.New()
	err := dbClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	grpcClient := grpcClient.New()
	err = grpcClient.Connect("bot.crypto-knight.site:50051")
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClient.Close()

	httpClient := httpClient.New()

	repo := repository.New(dbClient.DB)
	svc := service.New(repo, httpClient)

	svc.Kline.LoadKlinesForPeriod()
	log.Println("History is full")
	go func() {
		for {
			utilities.SleepUntilNextHour()
			svc.Kline.LoadKlinesForPeriod()
		}
	}()

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	// resp, err := grpcClient.History.ProcessHistory(ctx, &pbHistory.GetHistoryRequest{
	// 	Symbol: "BTCUSDT",
	// })
	// if err != nil {
	// 	log.Fatalf("request failed: %v", err)
	// }

	<-quit

	log.Println("History is stopped")
}
