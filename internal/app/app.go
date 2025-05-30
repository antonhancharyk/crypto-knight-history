package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/db"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service"
	grpcClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/grpc/client"
	httpClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/http/client"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/server"
	"github.com/antonhancharyk/crypto-knight-history/pkg/utilities"
	"github.com/joho/godotenv"
)

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbClient := db.New()
	err = dbClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	grpcClient := grpcClient.New()
	err = grpcClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClient.Close()

	httpClient := httpClient.New()

	repo := repository.New(dbClient.DB)
	svc := service.New(repo, httpClient, grpcClient)

	httpServer := server.New(svc)
	go func() {
		err := httpServer.Start()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	err = svc.Kline.LoadKlinesForPeriod()
	if err != nil {
		log.Print(err)
	}
	log.Println("initial klines are loaded")
	go func() {
		for {
			utilities.SleepUntilNextHour()
			err := svc.Kline.LoadKlinesForPeriod()
			if err != nil {
				log.Print(err)
			}
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
