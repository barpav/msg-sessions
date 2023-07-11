package users

import (
	"context"
	"fmt"
	"net"
	"time"

	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg  *Config
	conn *grpc.ClientConn
	stub usgrpc.UsersClient
}

func (c *Client) Connect() (err error) {
	c.cfg = &Config{}
	c.cfg.Read()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.conn, err = grpc.DialContext(ctx, net.JoinHostPort(c.cfg.host, c.cfg.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return fmt.Errorf("Can't connect to 'users' service: %s", err)
	}

	c.stub = usgrpc.NewUsersClient(c.conn)

	return nil
}

func (c *Client) Disconnect(ctx context.Context) (err error) {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
