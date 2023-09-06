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
	cfg  *config
	conn *grpc.ClientConn
	stub ssgrpc.SessionsClient
}

func (c *Client) Connect() (err error) {
	c.cfg = &config{}
	c.cfg.Read()

	target := net.JoinHostPort(c.cfg.host, c.cfg.port)
	optionCredentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	optionBlock := grpc.WithBlock()

	for try := 0; try < c.cfg.connTries; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.connTimeout)*time.Second)
		c.conn, err = grpc.DialContext(ctx, target, optionCredentials, optionBlock)
		cancel()

		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("can't connect to 'sessions' service: %w", err)
	}

	c.stub = ssgrpc.NewSessionsClient(c.conn)

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
		err = fmt.Errorf("failed to disconnect from 'sessions' service: %w", err)
	}

	return err
}
