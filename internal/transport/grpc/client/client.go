package client

import (
	"context"
	"time"

	"github.com/antonhancharyk/crypto-knight-history/internal/config"
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

func (c *Client) Connect(cfg config.GRPC) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Host,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*1024*1024),
			grpc.MaxCallSendMsgSize(50*1024*1024),
		),
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
		return c.conn.Close()
	}

	return nil
}
