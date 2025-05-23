package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	DB *sqlx.DB
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	log.Println("DB client is connecting...")

	connStr := fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable password=%s host=localhost",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	c.DB = db

	log.Println("DB client is connected")

	return nil
}

func (c *Client) Close() error {
	if c.DB != nil {
		log.Println("DB connection is closing...")
		return c.DB.Close()
	}

	return nil
}
