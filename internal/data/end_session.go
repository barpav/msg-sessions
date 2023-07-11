package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (s *Storage) EndSession(ctx context.Context, userId string, sessionId int64) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to end session: %w", err)
		}
	}()

	var l *lock
	l, err = s.lock(ctx, sessionsLockKey(userId), 30*time.Second)

	if err != nil {
		return err
	}

	defer func() {
		if err := l.unlock(); err != nil {
			log.Err(err).Msg("Failed to unlock user sessions.")
		}
	}()

	var keySum string
	keySum, err = s.db.HGet(ctx, sessionInfoKey(userId, sessionId), "keySum").Result()

	if err != nil && !errors.Is(err, redis.Nil) { // operation is idempotent
		return err
	}

	pipe := s.db.Pipeline()

	pipe.Del(ctx, keySum, sessionInfoKey(userId, sessionId))

	pipe.SRem(ctx, sessionsIdsKey(userId), sessionId)

	pipe.Decr(ctx, sessionsTotalKey(userId))

	commands, err := pipe.Exec(ctx)

	if err == nil {
		for _, cmd := range commands {
			err = errors.Join(err, cmd.Err())
		}
	}

	return err
}
