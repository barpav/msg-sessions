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

type NewSession struct {
	User string

	id int64
}

var ErrTooManySessions = errors.New("too many sessions")

func (m *NewSession) Create(ctx context.Context, s *Storage) (session *ActiveSession, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to create new session: %w", err)
		}
	}()

	if m.User == "" {
		return nil, errors.New("user must be specified")
	}

	var l *lock
	l, err = s.lock(ctx, sessionsLockKey(m.User), 30*time.Second)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := l.unlock(); err != nil {
			log.Err(err).Msg("Failed to unlock user sessions.")
		}
	}()

	var sessionsTotal int
	sessionsTotalKey := sessionsTotalKey(m.User)

	sessionsTotal, err = s.db.Get(ctx, sessionsTotalKey).Int()

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if sessionsTotal >= s.cfg.maxSessionsPerUser {
		return nil, ErrTooManySessions
	}

	var id int64
	id, err = s.db.Incr(ctx, sessionsLastIdKey(m.User)).Result()

	if err != nil {
		return nil, err
	}

	newSessionKey := uuid.NewString()
	newSessionKeySum := sessionKeySum(newSessionKey)

	var set bool
	set, err = s.db.HSetNX(ctx, newSessionKeySum, "user", m.User).Result()

	if err != nil {
		return nil, err
	}

	if !set {
		return nil, errors.New("new session key not unique")
	}

	pipe := s.db.Pipeline()

	pipe.HSet(ctx, newSessionKeySum, "id", id)

	pipe.SAdd(ctx, sessionsIdsKey(m.User), id)

	pipe.Incr(ctx, sessionsTotalKey)

	commands, err := pipe.Exec(ctx)

	if err == nil {
		for _, cmd := range commands {
			err = errors.Join(err, cmd.Err())
		}
	}

	if err != nil {
		return nil, err
	}

	session, m.id = &ActiveSession{Key: newSessionKey}, id

	return session, nil
}

func (m *NewSession) Id() int64 {
	return m.id
}
