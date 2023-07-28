package sessions

import (
	"context"
	"fmt"

	ssgrpc "github.com/barpav/msg-sessions/sessions_service_go_grpc"
)

func (c *Client) EndAllSessions(ctx context.Context, userId string) (err error) {
	_, err = c.stub.EndAll(ctx, &ssgrpc.User{Id: userId})

	if err != nil {
		return fmt.Errorf("failed to end all user '%s' sessions (service 'sessions' client): %w", userId, err)
	}

	return nil
}
