package db

import (
	"fmt"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	DB *sqlx.DB
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connect(cfg config.DB) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
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
		return c.DB.Close()
	}

	return nil
}
