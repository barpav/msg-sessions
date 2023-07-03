package data

import (
	"context"
	"fmt"
)

type ErrSessionNotFound struct{}

func (s *Storage) SessionKeyInfo(ctx context.Context, key string) (userId string, sessionId int64, err error) {
	defer func() {
		if err != nil {
			if _, ok := err.(*ErrSessionNotFound); ok {
				return
			}
			err = fmt.Errorf("Failed to get session key info: %w", err)
		}
	}()

	keySum := sessionKeySum(key)

	pipe := s.db.Pipeline()

	userIdCmd := s.db.HGet(ctx, keySum, "user")
	sessionIdCmd := s.db.HGet(ctx, keySum, "id")

	_, err = pipe.Exec(ctx)

	if err != nil {
		return "", 0, err
	}

	userId, err = userIdCmd.Result()

	if err != nil {
		return "", 0, err
	}

	sessionId, err = sessionIdCmd.Int64()

	if err != nil {
		return "", 0, err
	}

	if userId == "" || sessionId == 0 {
		return "", 0, &ErrSessionNotFound{}
	}

	return userId, sessionId, nil
}

func (e *ErrSessionNotFound) Error() string {
	return "session not found"
}

func (e *ErrSessionNotFound) ImplementsSessionNotFoundError() {
}
