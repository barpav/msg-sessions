package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type ErrTooManySessions struct{}

func (s *Storage) StartNewSession(ctx context.Context, userId, ip, agent string) (id int64, key string, err error) {
	defer func() {
		if err != nil {
			if _, ok := err.(*ErrTooManySessions); ok {
				return
			}
			err = fmt.Errorf("Failed to start new session: %w", err)
		}
	}()

	var l *lock
	l, err = s.lock(ctx, sessionsLockKey(userId), 30*time.Second)

	if err != nil {
		return 0, "", err
	}

	defer func() {
		if err := l.unlock(); err != nil {
			log.Err(err).Msg("Failed to unlock user sessions.")
		}
	}()

	var sessionsTotal int
	sessionsTotalKey := sessionsTotalKey(userId)

	sessionsTotal, err = s.db.Get(ctx, sessionsTotalKey).Int()

	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, "", err
	}

	if sessionsTotal >= s.cfg.maxSessionsPerUser {
		return 0, "", &ErrTooManySessions{}
	}

	id, err = s.db.Incr(ctx, sessionsLastIdKey(userId)).Result()

	if err != nil {
		return 0, "", err
	}

	key = uuid.NewString()
	keySum := sessionKeySum(key)

	var set bool
	set, err = s.db.HSetNX(ctx, keySum, "user", userId).Result()

	if err != nil {
		return 0, "", err
	}

	if !set {
		return 0, "", errors.New("new session key not unique")
	}

	pipe := s.db.Pipeline()

	pipe.HSet(ctx, keySum, "id", id)

	pipe.SAdd(ctx, sessionsIdsKey(userId), id)

	now := time.Now().UTC()
	s.db.HSet(ctx, sessionInfoKey(userId, id),
		"keySum", keySum,
		"created", now,
		"lastActivity", now,
		"lastIp", ip,
		"lastAgent", agent).Err()

	pipe.Incr(ctx, sessionsTotalKey)

	commands, err := pipe.Exec(ctx)

	if err == nil {
		for _, cmd := range commands {
			err = errors.Join(err, cmd.Err())
		}
	}

	if err != nil {
		return 0, "", err
	}

	return id, key, err
}

func (e *ErrTooManySessions) Error() string {
	return "maximum number of user sessions exceeded"
}

func (e *ErrTooManySessions) ImplementsTooManySessionsError() {
}
