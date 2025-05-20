package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/antonhancharyk/crypto-knight-history/internal/db"
	"github.com/joho/godotenv"
)

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	godotenv.Load()

	db := db.Connect()
	defer db.Close()

	// log.Println("History is running...")

	<-quit

	// log.Println("History is stopping...")
}
