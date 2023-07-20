package sessions

import (
	"context"
	"fmt"
	"net"
	"time"

	ssgrpc "github.com/barpav/msg-sessions/sessions_service_go_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg  *Config
	conn *grpc.ClientConn
	stub ssgrpc.SessionsClient
}

func (c *Client) Connect() (err error) {
	c.cfg = &Config{}
	c.cfg.Read()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.conn, err = grpc.DialContext(ctx, net.JoinHostPort(c.cfg.host, c.cfg.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return fmt.Errorf("Can't connect to 'sessions' service: %w", err)
	}

	c.stub = ssgrpc.NewSessionsClient(c.conn)

	return nil
}

func (c *Client) Disconnect(ctx context.Context) (err error) {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
