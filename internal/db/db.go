package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() *sqlx.DB {
	log.Println("DB is running...")

	connectStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=localhost", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	log.Println("DB is run")

	return db
}
