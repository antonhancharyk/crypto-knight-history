package client

import (
	"context"
	"log"
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

func (c *Client) Connect(addr string) error {
	log.Println("gRPC client is connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.History = pbHistory.NewHistoryServiceClient(conn)

	log.Println("gRPC client is connected")

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		log.Println("gRPC client is closing connection...")
		return c.conn.Close()
	}

	return nil
}
