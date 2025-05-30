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
		"user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
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

	return nil
}

func (c *Client) Close() error {
	if c.DB != nil {
		log.Println("shutting down DB client connection...")
		return c.DB.Close()
	}

	return nil
}
