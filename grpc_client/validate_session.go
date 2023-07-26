package sessions

import (
	"context"
	"fmt"

	ssgrpc "github.com/barpav/msg-sessions/sessions_service_go_grpc"
)

func (c *Client) ValidateSession(ctx context.Context, key, ip, agent string) (userId string, err error) {
	var result *ssgrpc.ValidationResult
	result, err = c.stub.Validate(ctx, &ssgrpc.SessionData{Key: key, Ip: ip, Agent: agent})

	if err != nil {
		return "", fmt.Errorf("failed to validate session (service 'sessions' client): %w", err)
	}

	return result.User, nil
}
