package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/client/binance"
	"github.com/antonhancharyk/crypto-knight-history/internal/config"
	"github.com/antonhancharyk/crypto-knight-history/internal/constant"
	"github.com/antonhancharyk/crypto-knight-history/internal/db"
	"github.com/antonhancharyk/crypto-knight-history/internal/repository"
	"github.com/antonhancharyk/crypto-knight-history/internal/service"
	grpcClient "github.com/antonhancharyk/crypto-knight-history/internal/transport/grpc/client"
	"github.com/antonhancharyk/crypto-knight-history/internal/transport/http/server"
	"github.com/joho/godotenv"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	if err := godotenv.Load(); err != nil {
		log.Print("optional .env load: ", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbClient := db.New()
	if err := dbClient.Connect(cfg.DB); err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	grpcCli := grpcClient.New()
	if err := grpcCli.Connect(cfg.GRPC); err != nil {
		log.Fatal(err)
	}
	defer grpcCli.Close()

	binanceClient := binance.New()
	repo := repository.New(dbClient.DB)
	svc := service.New(repo.Kline, binanceClient, grpcCli.History)

	httpServer := server.New(svc, cfg.Server)
	go func() {
		err := httpServer.Start()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	if false {
		startAll := time.Now()
		for _, interval := range constant.KLINE_INTERVALS {
			start := time.Now()

			log.Printf("%s started at %s", interval, start.UTC().Format(time.RFC3339))

			err := svc.Kline.LoadInterval(interval)
			if err != nil {
				log.Print(err)
				continue
			}

			finish := time.Now()
			log.Printf(
				"%s finished at %s (duration: %s)",
				interval,
				finish.UTC().Format(time.RFC3339),
				finish.Sub(start).Round(time.Millisecond),
			)
		}
		log.Printf("all klines loaded in %s", time.Since(startAll).Round(time.Millisecond))
	}

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
