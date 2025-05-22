package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/antonhancharyk/crypto-knight-history/internal/db"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
	"github.com/antonhancharyk/crypto-knight-history/pkg/utilities"
	"github.com/joho/godotenv"
)

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("History is running...")

	godotenv.Load("../.env")

	db := db.Connect()
	defer db.Close()

	repo := repository.New(db)
	httpClient := client.New()

	svc := service.New(repo, httpClient)

	log.Println("History is run")

	svc.Kline.LoadKlinesForPeriod()
	log.Println("History is full")
	go func() {
		for {
			utilities.SleepUntilNextHour()
			svc.Kline.LoadKlinesForPeriod()
		}
	}()

	<-quit

	log.Println("History is stopping...")
	log.Println("History is stopped")
}
