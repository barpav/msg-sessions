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
	cfg  *config
	conn *grpc.ClientConn
	stub usgrpc.UsersClient
}

func (c *Client) Connect() (err error) {
	c.cfg = &config{}
	c.cfg.Read()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.connTimeout)*time.Second)
	defer cancel()

	c.conn, err = grpc.DialContext(ctx, net.JoinHostPort(c.cfg.host, c.cfg.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return fmt.Errorf("can't connect to 'users' service: %w", err)
	}

	c.stub = usgrpc.NewUsersClient(c.conn)

	return nil
}

func (c *Client) Disconnect(ctx context.Context) (err error) {
	if c.conn == nil {
		return nil
	}

	var dErr error
	closed := make(chan struct{}, 1)

	go func() {
		dErr = c.conn.Close()
		closed <- struct{}{}
	}()

	select {
	case <-closed:
		err = dErr
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		err = fmt.Errorf("failed to disconnect from 'users' service: %w", err)
	}

	return err
}
