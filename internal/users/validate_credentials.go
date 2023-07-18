package users

import (
	"context"
	"fmt"

	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
)

func (c *Client) ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error) {
	result, err := c.stub.Validate(ctx, &usgrpc.Credentials{Id: userId, Password: password})

	if err != nil {
		return false, fmt.Errorf("Failed to validate credentials (service 'users' client): %w", err)
	}

	return result.Valid, nil
}
