package client

import (
	"context"
	"log"
	"os"
	"time"

	pbHistory "github.com/antongoncharik/crypto-knight-protos/gen/go/history"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	History pbHistory.HistoryServiceClient
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	log.Println("gRPC client is connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, os.Getenv("GRPC_HOST"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.History = pbHistory.NewHistoryServiceClient(conn)

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		log.Println("shutting down gRPC client connection...")
		return c.conn.Close()
	}

	return nil
}
